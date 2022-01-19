# open-esa

A web application that redirects human-readable URLs to [esa.io](https://esa.io/) posts.

## Getting Started

1. Register aplication. see the [documentation](https://docs.esa.io/posts/102#%E3%82%A2%E3%83%97%E3%83%AA%E3%82%B1%E3%83%BC%E3%82%B7%E3%83%A7%E3%83%B3%E3%81%AE%E7%99%BB%E9%8C%B2)
2. Write `open-esa.toml`
    ```sh
    cp open-esa.toml.esample open-esa.toml
    vi open-esa.toml
    ```
3. Run `open-esa`
4. Open http://localhost:8080/foo/bar/zoo

## Redirect Rules

* `http://your-open-esa.example.com/foo/bar/zoo`
    * Post exists:
    * `https://[team].esa.io/posts/[post num]`
    * Post does not exist:
    * `https://[team].esa.io/posts/posts/new?category_path=%2Ffoo%2Fbar%2Fzoo`
* `http://your-open-esa.example.com/foo/`
    * `https://[team].esa.io/#path=%2Ffoo`
