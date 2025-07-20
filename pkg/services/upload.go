package services

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/tgdrive/teldrive/internal/api"
	"github.com/tgdrive/teldrive/internal/auth"
	"github.com/tgdrive/teldrive/internal/crypt"
	"github.com/tgdrive/teldrive/internal/logging"
	"github.com/tgdrive/teldrive/internal/pool"
	"github.com/tgdrive/teldrive/internal/tgc"
	"go.uber.org/zap"

	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/message"
	"github.com/gotd/td/telegram/uploader"
	"github.com/gotd/td/tg"
	"github.com/tgdrive/teldrive/pkg/mapper"
	"github.com/tgdrive/teldrive/pkg/models"
)

var (
	saltLength      = 32
	ErrUploadFailed = errors.New("upload failed")
)

func (a *apiService) UploadsDelete(ctx context.Context, params api.UploadsDeleteParams) error {
	if err := a.db.Where("upload_id = ?", params.ID).Delete(&models.Upload{}).Error; err != nil {
		return &api.ErrorStatusCode{StatusCode: 500, Response: api.Error{Message: err.Error(), Code: 500}}
	}
	return nil
}

func (a *apiService) UploadsPartsById(ctx context.Context, params api.UploadsPartsByIdParams) ([]api.UploadPart, error) {
	parts := []models.Upload{}
	if err := a.db.Model(&models.Upload{}).Order("part_no").Where("upload_id = ?", params.ID).
		Where("created_at < ?", time.Now().UTC().Add(a.cnf.TG.Uploads.Retention)).
		Find(&parts).Error; err != nil {
		return nil, &apiError{err: err}
	}
	return mapper.ToUploadOut(parts), nil
}

func (a *apiService) UploadsStats(ctx context.Context, params api.UploadsStatsParams) ([]api.UploadStats, error) {
	userId := auth.GetUser(ctx)
	var stats []api.UploadStats
	err := a.db.Raw(`
    SELECT 
    dates.upload_date::date AS upload_date,
    COALESCE(SUM(files.size), 0)::bigint AS total_uploaded
    FROM 
        generate_series(
            (CURRENT_TIMESTAMP AT TIME ZONE 'UTC')::date - INTERVAL '1 day' * @days,
            (CURRENT_TIMESTAMP AT TIME ZONE 'UTC')::date,
            '1 day'
        ) AS dates(upload_date)
    LEFT JOIN 
    teldrive.files AS files
    ON 
        dates.upload_date = DATE_TRUNC('day', files.created_at)::date
        AND files.type = 'file'
        AND files.user_id = @userId
    GROUP BY 
        dates.upload_date
    ORDER BY 
        dates.upload_date
  `, sql.Named("days", params.Days-1), sql.Named("userId", userId)).Scan(&stats).Error

	if err != nil {
		return nil, &apiError{err: err}

	}
	return stats, nil
}

func (a *apiService) UploadsUpload(ctx context.Context, req *api.UploadsUploadReqWithContentType, params api.UploadsUploadParams) (*api.UploadPart, error) {
	var (
		channelId   int64
		err         error
		client      *telegram.Client
		token       string
		index       int
		channelUser string
		out         api.UploadPart
	)

	if params.Encrypted.Value && a.cnf.TG.Uploads.EncryptionKey == "" {
		return nil, &apiError{err: errors.New("encryption is not enabled"), code: 400}
	}

	userId := auth.GetUser(ctx)

	fileStream := req.Content.Data

	fileSize := params.ContentLength

	if params.ChannelId.Value == 0 {
		channelId, err = getDefaultChannel(a.db, a.cache, userId)
		if err != nil {
			return nil, err
		}
	} else {
		channelId = params.ChannelId.Value
	}

	tokens, err := getBotsToken(a.db, a.cache, userId, channelId)

	if err != nil {
		return nil, err
	}

	if len(tokens) == 0 {
		client, err = tgc.AuthClient(ctx, &a.cnf.TG, auth.GetJWTUser(ctx).TgSession)
		if err != nil {
			return nil, err
		}
		channelUser = strconv.FormatInt(userId, 10)
	} else {
		a.worker.Set(tokens, channelId)
		token, index = a.worker.Next(channelId)
		client, err = tgc.BotClient(ctx, a.tgdb, &a.cnf.TG, token)

		if err != nil {
			return nil, err
		}

		channelUser = strings.Split(token, ":")[0]
	}

	middlewares := tgc.NewMiddleware(&a.cnf.TG, tgc.WithFloodWait(),
		tgc.WithRecovery(ctx),
		tgc.WithRetry(a.cnf.TG.Uploads.MaxRetries),
		tgc.WithRateLimit())

	uploadPool := pool.NewPool(client, int64(a.cnf.TG.PoolSize), middlewares...)

	defer uploadPool.Close()

	logger := logging.FromContext(ctx)

	logger.Debug("uploading chunk",
		zap.String("fileName", params.FileName),
		zap.String("partName", params.PartName),
		zap.String("bot", channelUser),
		zap.Int("botNo", index),
		zap.Int("chunkNo", params.PartNo),
		zap.Int64("partSize", fileSize),
	)

	err = tgc.RunWithAuth(ctx, client, token, func(ctx context.Context) error {

		logger.Debug("starting upload process",
			zap.Int64("channelId", channelId),
			zap.String("partName", params.PartName))

		channel, err := tgc.GetChannelById(ctx, client.API(), channelId)

		if err != nil {
			logger.Error("failed to get channel",
				zap.Int64("channelId", channelId),
				zap.Error(err))
			return fmt.Errorf("failed to get channel %d: %w", channelId, err)
		}

		logger.Debug("channel retrieved successfully",
			zap.Int64("channelId", channel.ChannelID))

		var salt string

		if params.Encrypted.Value {
			logger.Debug("encryption enabled, generating salt")
			salt, _ = generateRandomSalt()
			cipher, err := crypt.NewCipher(a.cnf.TG.Uploads.EncryptionKey, salt)
			if err != nil {
				logger.Error("failed to create cipher", zap.Error(err))
				return fmt.Errorf("failed to create cipher: %w", err)
			}
			fileSize = crypt.EncryptedSize(fileSize)
			fileStream, err = cipher.EncryptData(fileStream)
			if err != nil {
				logger.Error("failed to encrypt data", zap.Error(err))
				return fmt.Errorf("failed to encrypt data: %w", err)
			}
			logger.Debug("data encrypted successfully", zap.Int64("encryptedSize", fileSize))
		}

		client := uploadPool.Default(ctx)

		logger.Debug("starting telegram upload",
			zap.Int64("fileSize", fileSize),
			zap.Int("threads", a.cnf.TG.Uploads.Threads))

		u := uploader.NewUploader(client).WithThreads(a.cnf.TG.Uploads.Threads).WithPartSize(512 * 1024)

		upload, err := u.Upload(ctx, uploader.NewUpload(params.PartName, fileStream, fileSize))

		if err != nil {
			logger.Error("telegram uploader failed",
				zap.Error(err),
				zap.String("errorType", fmt.Sprintf("%T", err)))
			return fmt.Errorf("telegram upload failed: %w", err)
		}

		logger.Debug("telegram upload completed successfully")

		document := message.UploadedDocument(upload).Filename(params.PartName).ForceFile(true)

		sender := message.NewSender(client)

		target := sender.To(&tg.InputPeerChannel{ChannelID: channel.ChannelID,
			AccessHash: channel.AccessHash})

		logger.Debug("sending media to channel",
			zap.Int64("channelId", channel.ChannelID),
			zap.String("filename", params.PartName))

		res, err := target.Media(ctx, document)

		if err != nil {
			logger.Error("failed to send media to channel",
				zap.Error(err),
				zap.String("errorType", fmt.Sprintf("%T", err)),
				zap.Int64("channelId", channel.ChannelID))
			return fmt.Errorf("failed to send media to channel: %w", err)
		}

		logger.Debug("media sent successfully, processing response")

		updates, ok := res.(*tg.Updates)
		if !ok {
			logger.Error("unexpected response type from telegram",
				zap.String("responseType", fmt.Sprintf("%T", res)))
			return fmt.Errorf("unexpected response type from telegram upload: %T", res)
		}

		logger.Debug("received updates response",
			zap.Int("updatesCount", len(updates.Updates)))

		var message *tg.Message

		for i, update := range updates.Updates {
			logger.Debug("processing update",
				zap.Int("updateIndex", i),
				zap.String("updateType", fmt.Sprintf("%T", update)))
			
			channelMsg, ok := update.(*tg.UpdateNewChannelMessage)
			if ok {
				message = channelMsg.Message.(*tg.Message)
				logger.Debug("found channel message",
					zap.Int("messageId", message.ID))
				break
			}
		}

		if message == nil || message.ID == 0 {
			logger.Error("no valid message received from telegram",
				zap.Bool("messageIsNil", message == nil),
				zap.Int("messageId", func() int { if message != nil { return message.ID } else { return 0 } }()))
			return fmt.Errorf("upload failed: no valid message received")
		}

		logger.Debug("message created successfully",
			zap.Int("messageId", message.ID))

		partUpload := &models.Upload{
			Name:      params.PartName,
			UploadId:  params.ID,
			PartId:    message.ID,
			ChannelId: channelId,
			Size:      fileSize,
			PartNo:    int(params.PartNo),
			UserId:    userId,
			Encrypted: params.Encrypted.Value,
			Salt:      salt,
		}

		logger.Debug("saving upload record to database",
			zap.Int("partId", message.ID),
			zap.Int64("size", fileSize))

		if err := a.db.Create(partUpload).Error; err != nil {
			logger.Error("database insert failed, cleaning up uploaded message",
				zap.Error(err),
				zap.Int("messageId", message.ID))
			// Clean up the uploaded message if database insert fails
			client.ChannelsDeleteMessages(ctx, &tg.ChannelsDeleteMessagesRequest{Channel: channel, ID: []int{message.ID}})
			return fmt.Errorf("database insert failed: %w", err)
		}

		logger.Debug("upload record saved successfully, verifying upload")

		// Verify the upload by fetching the message back
		v, err := client.ChannelsGetMessages(ctx, &tg.ChannelsGetMessagesRequest{Channel: channel, ID: []tg.InputMessageClass{&tg.InputMessageID{ID: message.ID}}})

		if err != nil || v == nil {
			logger.Error("failed to verify upload by fetching message",
				zap.Error(err),
				zap.Bool("responseIsNil", v == nil),
				zap.Int("messageId", message.ID))
			return ErrUploadFailed
		}

		logger.Debug("verification request completed, checking response")

		switch msgs := v.(type) {
		case *tg.MessagesChannelMessages:
			if len(msgs.Messages) == 0 {
				logger.Error("verification failed: no messages returned",
					zap.Int("messageId", message.ID))
				return ErrUploadFailed
			}
			
			logger.Debug("verification successful, checking document",
				zap.Int("messagesCount", len(msgs.Messages)))
			
			doc, ok := msgDocument(msgs.Messages[0])
			if !ok {
				logger.Error("verification failed: message does not contain valid document",
					zap.String("messageType", fmt.Sprintf("%T", msgs.Messages[0])))
				return ErrUploadFailed
			}
			
			logger.Debug("document found in verification",
				zap.Int64("documentSize", doc.Size),
				zap.Int64("expectedSize", fileSize))
			
			if doc.Size != fileSize {
				logger.Error("verification failed: size mismatch, cleaning up",
					zap.Int64("documentSize", doc.Size),
					zap.Int64("expectedSize", fileSize),
					zap.Int("messageId", message.ID))
				client.ChannelsDeleteMessages(ctx, &tg.ChannelsDeleteMessagesRequest{Channel: channel, ID: []int{message.ID}})
				return ErrUploadFailed
			}
			
			logger.Debug("verification passed: sizes match")
		default:
			logger.Error("verification failed: unexpected message type",
				zap.String("messageType", fmt.Sprintf("%T", v)))
			return ErrUploadFailed
		}

		out = api.UploadPart{
			Name:      partUpload.Name,
			PartId:    partUpload.PartId,
			ChannelId: partUpload.ChannelId,
			PartNo:    partUpload.PartNo,
			Size:      partUpload.Size,
			Encrypted: partUpload.Encrypted,
		}
		out.SetSalt(api.NewOptString(partUpload.Salt))
		
		logger.Debug("upload process completed successfully",
			zap.Int("partId", partUpload.PartId),
			zap.Int64("size", partUpload.Size))
		
		return nil

	})

	if err != nil {
		logger.Error("upload failed", 
			zap.String("fileName", params.FileName),
			zap.String("partName", params.PartName),
			zap.Int("chunkNo", params.PartNo),
			zap.Error(err),
			zap.String("errorType", fmt.Sprintf("%T", err)),
			zap.Int64("channelId", channelId),
			zap.String("bot", channelUser),
			zap.Int("botIndex", index))
		
		// Return detailed error information for frontend handling
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, &apiError{err: errors.New("upload timeout - please retry with smaller chunk size"), code: 408}
		}
		
		if strings.Contains(err.Error(), "FLOOD_WAIT") {
			return nil, &apiError{err: errors.New("telegram rate limit exceeded - please wait and retry"), code: 429}
		}
		
		if strings.Contains(err.Error(), "FILE_PART_") {
			return nil, &apiError{err: errors.New("invalid file part - please retry upload"), code: 400}
		}
		
		if strings.Contains(err.Error(), "CHANNEL_PRIVATE") {
			return nil, &apiError{err: errors.New("channel access denied - check bot permissions"), code: 403}
		}
		
		return nil, &apiError{err: fmt.Errorf("upload failed: %w", err), code: 500}
	}
	logger.Debug("upload finished", zap.String("fileName", params.FileName),
		zap.String("partName", params.PartName),
		zap.Int("chunkNo", params.PartNo))
	return &out, nil

}

func msgDocument(m tg.MessageClass) (*tg.Document, bool) {
	res, ok := m.AsNotEmpty()
	if !ok {
		return nil, false
	}
	msg, ok := res.(*tg.Message)
	if !ok {
		return nil, false
	}

	media, ok := msg.Media.(*tg.MessageMediaDocument)
	if !ok {
		return nil, false
	}
	return media.Document.AsNotEmpty()
}

func generateRandomSalt() (string, error) {
	randomBytes := make([]byte, saltLength)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	hasher := sha256.New()
	hasher.Write(randomBytes)
	hashedSalt := base64.URLEncoding.EncodeToString(hasher.Sum(nil))

	return hashedSalt, nil
}
