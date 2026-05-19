# Integrations

`integrations/lcp_forwarder.py` is a small dependency-free bridge for tools that can call a script or webhook after a book import.

```bash
python3 integrations/lcp_forwarder.py "/library/book.epub" --title "Book Title"
```

Use it from Calibre post-import automation, calibre-web hooks, Kavita import automation, or a tiny webhook wrapper of your own.

| Variable | Default |
| --- | --- |
| `LCP_BASE_URL` | `http://localhost:8080` |
| `LCP_USERNAME` | `publisher` |
| `LCP_PASSWORD` | `publisher` |
| `LCP_2FA_CODE` | empty |
