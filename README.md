# Teldrive
[![Ask DeepWiki](https://deepwiki.com/badge.svg)](https://deepwiki.com/tgdrive/teldrive)

Teldrive is a powerful utility that enables you to organise your telegram files and much more.

## Advantages Over Alternative Solutions

- **Exceptional Speed:** Teldrive stands out among similar tools, thanks to its implementation in Go, a language known for its efficiency. Its performance surpasses alternatives written in Python and other languages, with the exception of Rust.

- **Enhanced Management Capabilities:** Teldrive not only excels in speed but also offers an intuitive user interface for efficient file interaction which other tool lacks. Its compatibility with Rclone further enhances file management.

> [!IMPORTANT]
> Teldrive functions as a wrapper over your Telegram account, simplifying file access. However, users must adhere to the limitations imposed by the Telegram API. Teldrive is not responsible for any consequences arising from non-compliance with these API limits.You will be banned instantly if you misuse telegram API.

## Quick Start

### Production Deployment
```bash
git clone https://github.com/tgdrive/teldrive.git
cd teldrive
docker-compose up -d
```

### Development Environment
```bash
git clone https://github.com/tgdrive/teldrive.git
cd teldrive
docker-compose -f docker-compose.dev.yml up
```

Visit https://teldrive-docs.pages.dev for detailed setup instructions.

# Recognitions

<a href="https://trendshift.io/repositories/7568" target="_blank"><img src="https://trendshift.io/api/badge/repositories/7568" alt="divyam234%2Fteldrive | Trendshift" style="width: 250px; height: 55px;" width="250" height="55"/></a>

## Best Practices for Using Teldrive

### Dos:

- **Follow Limits:** Adhere to the limits imposed by Telegram servers to avoid account bans and automatic deletion of your channel.Your files will be removed from telegram servers if you try to abuse the service as most people have zero brains they will still do so good luck.
- **Responsible Storage:** Be mindful of the content you store on Telegram. Utilize storage efficiently and only keep data that serves a purpose.
  
### Don'ts:
- **Data Hoarding:** Avoid excessive data hoarding, as it not only violates Telegram's terms.
  
By following these guidelines, you contribute to the responsible and effective use of Telegram, maintaining a fair and equitable environment for all users.

## Contributing

Feel free to contribute to this project.See [CONTRIBUTING.md](CONTRIBUTING.md) for more information.

## Donate

If you like this project small contribution would be appreciated [Paypal](https://paypal.me/redux234).

## Star History

<a href="https://www.star-history.com/#tgdrive/teldrive&Date">
 <picture>
   <source media="(prefers-color-scheme: dark)" srcset="https://api.star-history.com/svg?repos=tgdrive/teldrive&type=Date&theme=dark" />
   <source media="(prefers-color-scheme: light)" srcset="https://api.star-history.com/svg?repos=tgdrive/teldrive&type=Date" />
   <img alt="Star History Chart" src="https://api.star-history.com/svg?repos=tgdrive/teldrive&type=Date" />
 </picture>
</a>
