# oblig-1



## Info documentation

### Structure

This path retrieves all cities in the country with the given iso-2 code.

The info page provides the following fields in **`json`** format for a given country: 
+ name (string)
+ continents (list of strings)
+ languages (object)
+ population (integer)
+ borders (list of strings)
+ flag (string)
+ capital (list of strings)
+ cities (list of strings)

### Syntax

The general syntax is `http://localhost:8080/countryinfo/v1/info/{:iso}`               ***CHANGE LINK***  

`{:iso}` is the two-letter iso code for a given country, eg. `no` for Norway.

Optionally, you can set a limit for how many cities to be retrieved. When no limit is given, it will be 10 by default.

The syntax for setting a limit is `http://localhost:8080/countryinfo/v1/info/{:iso}{?limit={:limit}}` where the whole `{?limit=...}` block is optional, and `{:limit}` is the limit, eg. `30`.

Leaving `{?limit}` empty as in `...?limit` or `...?limit=` will be treated as having no limit, and will return all cities.

### Example 

The call to

`http://localhost:8080/countryinfo/v1/info/no?limit=25`

will give the following **`json`**:

```json
{
    "name": "Norway",
    "continents": [
        "Europe"
    ],
    "languages": {
        "nno": "Norwegian Nynorsk",
        "nob": "Norwegian Bokmål",
        "smi": "Sami"
    },
    "population": 5379475,
    "borders": [
        "FIN",
        "SWE",
        "RUS"
    ],
    "flag": "🇳🇴",
    "capital": [
        "Oslo"
    ],
    "cities": [
        "Abelvaer",
        "Adalsbruk",
        "Adland",
        "Agotnes",
        "Agskardet",
        "Aker",
        "Akkarfjord",
        "Akrehamn",
        "Al",
        "Alen",
        "Algard",
        "Almas",
        "Alta",
        "Alvdal",
        "Amli",
        "Amot",
        "Ana-Sira",
        "Andalsnes",
        "Andenes",
        "Angvika",
        "Ankenes",
        "Annstad",
        "Ardal",
        "Ardalstangen",
        "Arendal"
    ]
}
```





## Population documentation

This path retrives the population count of the country with the given iso-2 code for each recorded year, and the average population count in this timeframe.

The population page povides the following in **`json`** format:
+ mean (integer)
+ data (object)
    + populationCounts (list of strings)


### Syntax

The general syntax is `http://localhost:8080/countryinfo/v1/population/{:iso}` ***CHANGE LINK***  

`{:iso}` is the two-letter iso code for a given country, eg. `no` for Norway.

Optionally you can set a bound for the start and end year.  
This limit will include both the start and end year in the retrived output.  
This means that if you provide eg. '`1990-2003`', all years in that span including 1990 and 2003 will be retrived.  
If the provided minimum is lower than the first recorded year, the first recorded year will be the minimum bound.  

When there is no provided limit, the bounds will be the earliest recorded year as the minimum bound and the current year as the maximum bound.

The syntax for setting a limit is `http://localhost:8080/countryinfo/v1/population/{:iso}{?limit={:startYear-endYear}}` ***CHANGE LINK***  
The `{?limit=...}` block is optional, and `{:startYear-endYear}` is the start and end boundaries.  
***NOTE:*** the '`-`' is mandatory if a limit is provided.


Leaving `{?limit}` empty as in `...?limit` or `...?limit=` will be treated as if no limit was provided, and the boundaries will be the earliest recorded year and the current year.



### Example

The call to  

`http://localhost:8080/countryinfo/v1/population/no?limit=1999-2004`

will give the following **`json`**:

```json
{
    "mean": 4526925,
    "data": {
        "populationCounts": [
            {
                "year": 1999,
                "value": 4461913
            },
            {
                "year": 2000,
                "value": 4490967
            },
            {
                "year": 2001,
                "value": 4513751
            },
            {
                "year": 2002,
                "value": 4538159
            },
            {
                "year": 2003,
                "value": 4564855
            },
            {
                "year": 2004,
                "value": 4591910
            }
        ]
    }
}
```





## Error documentation


### Error codes


#### Arising from '/population/' path

##### 200

The provided input in the `{iso}` field is not a two letter iso-2 code.  
Make sure it is two letters long.

Example of a valid input:
`http://localhost:8080/countryinfo/v1/population/gb` ***CHANGE LINK***



##### 201

The provided input int the `{iso}` field is valid syntactically, but it didn't match a valid iso-2 code.  
For a list of valid iso-2 codes, see **[THIS](https://en.wikipedia.org/wiki/ISO_3166-1_alpha-2#Officially_assigned_code_elements)** link.



##### 202

Two arguments was expected in `{:startYear-endYear}`, but some other amount was provided.  
Make sure to provide two numbers, separated by a '`-`'  

Example scenarios where this error might occur:  
> `?limit=2000` - too few arguments (1)
>
> `?limit=1997-2003-2009` - too many arguments (3)



##### 203

One or both arguments in `{:startYear-endYear}` are empty.  
Make sure to provide a number before and after the '`-`'.  
It does not matter what the number is, as long as it is between 0 and the 64-bit integer upper bound (9223372036854775807), as anything above that will not be recognised as a number.  

Example scenarios where this error might occur:  
> `?limit=2003-` - no end year provided
> 
> `?limit=-2003` - no start year provided
>
> `?limit=-` - neither start year nor end year provided



##### 204.1

The first argument of `{:startYear-endYear}` - `startYear` - is not a number.  
Whitespaces and/or letters are not permitted.  

Example scenarios where this error might occur:  
> `?limit=fish-2004` - fish is not a number
> 
> `?limit=2009 -2016` - whitespaces are disallowed
>
> `?limit=twothousandandthree-3500` - spelling out numbers does not work unfortunately



##### 204.2

The second argument of `{:startYear-endYear}` - `endYear` - is not a number.  
Whitespaces and/or letters are not permitted.

Example scenarios where this error might occur:  
> `?limit=0-Christopher Franz` - Christopher Franz is not a number
> 
> `?limit=2004- 2023` - whitespaces are disallowed
>
> `?limit=1995-9223372036854775808` - 9223372036854775808 is outside the 64-bit integer upper limit (9223372036854775807)



##### 205

`startYear` is greater than `endYear` in `{:startYear-endYear}`.  
Start year must be smaller for a valid range to be possible.

Example scenarios where this error might occur:  
> `?limit=2023-1999` - 2023 is greater than 1999
> 
> `?limit=5000-2000` - 5000 is greater than 2000



