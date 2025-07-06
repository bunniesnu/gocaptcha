# GoCaptcha

<a href="https://github.com/bunniesnu/gocaptcha/tags"><img src="https://img.shields.io/github/tag/bunniesnu/gocaptcha.svg" alt="Latest Tag"></a>
<a href="https://github.com/bunniesnu/gocaptcha/actions"><img src="https://github.com/bunniesnu/gocaptcha/actions/workflows/test-schedule.yml/badge.svg" alt="Test Status"></a>

A Go port of [PyPasser](https://github.com/xHossein/PyPasser/tree/master).

Bypass reCaptchaV3 only by calling a function. This does not support reCaptchaV2.

## Installation

```
go get github.com/bunniesnu/gocaptcha
```

## Usage

```
recaptcha, err := NewRecaptchaV3("your_anchor_url", proxy, t * time.Second)
token, err := recaptcha.Solve()
```

Proxy is not test at this point. Please provide ```nil```.

## Legal Disclaimer

This was made for educational purposes only, nobody which directly involved in this project is responsible for any damages caused.
**You are responsible for your actions.**

## License

[MIT](https://choosealicense.com/licenses/mit/)