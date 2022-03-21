# esaop

A web application that redirects human-readable URLs to [esa.io](https://esa.io/) posts.

## Getting Started

1. Register your aplication. see the [documentation](https://docs.esa.io/posts/102#%E3%82%A2%E3%83%97%E3%83%AA%E3%82%B1%E3%83%BC%E3%82%B7%E3%83%A7%E3%83%B3%E3%81%AE%E7%99%BB%E9%8C%B2)
    * Redirect URL: `http://[team].esa.io/auth/callback`
3. Write `esaop.toml`
    ```sh
    cp esaop.toml.sample esaop.toml
    vi esaop.toml
    ```
3. Run `esaop`
4. Open http://localhost:8080/foo/bar/zoo

## Redirect Rules

* `http://your-esaop.example.com/foo/bar/zoo`
    * Post exists:
    * -> `https://[team].esa.io/posts/[post num]`
    * Post does not exist:
    * -> `https://[team].esa.io/posts/posts/new?category_path=%2Ffoo%2Fbar%2Fzoo`
* `http://your-esaop.example.com/foo/`
    * -> `https://[team].esa.io/#path=%2Ffoo`

### Use Date/CRON expression

Access date: `2022/03/21`

* `http://[team].esa.io/${yyyy/MM/dd}`
  * -> Path: `2022/03/21`
* `http://[team].esa.io/${*,*m*,*,5|yyyy/MM/dd}`
  * -> `Path: 2022/03/25`
* `http://[team].esa.io/${*,*,10,*,*|yyyy/MM/dd}`
  * -> `Path: 2022/04/10`
