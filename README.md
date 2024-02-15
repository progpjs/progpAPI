# ProgpJS ProgpAPI

It's the low level API for the javascript engine. It's mainly used when you want to add special features
to a Go function exposed to the javascript engine, like progpAPI.StringBuffer which allows to return a
buffer when a string is expected (it's allow avoiding buffer to string conversion).

The split between ProgpAPI and progpScripts has this logic: ProgpJS modules must only access progpAPI
and never progpScripts. Mainly because an update to ProgpAPI involves regenerating all the pre-compiled
Go plugins, theses being very usefully since libProgpV8 take about 5 secondes to compiles without that.  

See https://github.com/progpjs/documentation for more information.