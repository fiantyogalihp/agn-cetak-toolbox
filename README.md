# agn-cetak-toolbox (Cetak Bukti Bayar)
is a web platform to replace wrong JSON to right JSON, this common issue on angon support.

### Format screen: 
Standard Format:

```json
{
  "screen_name": "<your_screen_name>",
  "arrange": [
    "1st_json_field",
    "2nd_json_field",
    "3rd_json_field",
    "4th_json_field",
    "etc"
  ],
  "required": [
    "<your_required_json_field_to_adjust>:<with_location>"
    "<your_index_location_data_in_array>:<your_field_want_to_be_check>", // if the value of the field is inside of array
    "etc",
  ],
  "adjustment": {
    "<your_destination>": "<your_source>",
    "etc": "etc"
  }
}
```

example:

```json
{
  "screen_name": "PBB Kab. Banjar",
  "arrange": [
    "inq",
    "amount",
    "refnum",
    "etc"
  ],
  "required": [
    "refnum",
    "inq:data",
    "inq:datetime",
    "0,0,2,1,0:PBB Th", // if the value of the field is inside of array
    "etc",
  ],
  "adjustment": {
    "pay:data": "inq:data",
    "etc": "etc"
  }
}
```
