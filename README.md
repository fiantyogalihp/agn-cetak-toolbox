# dynamic-json-parsing-struct
cetak bukti bayar depend on angon issue

## TODO
1. "/" load "Name" from all screen.json filename and replace it to the radio button htmx as "format validation"
2. when choose a radio button, its gonna be load the screen.json file and filename for explicit
3. send the format validation data inside the screen.json to backend which gonna be processed
4. format screen:


```json
{
  "arrange": [
    "datetime",
    "rp_tag",
    "dll ...",
  ],
  "adjustment": {
    "inq:receipt": "pay:receipt:0,0,2,1",
    "source": "dest"
  }
}
```
