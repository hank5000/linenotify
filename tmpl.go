package main

const authTmpl = `
<!DOCTYPE html>
<html lang="tw">
    <head>
        <title></title>
        <meta charset="UTF-8">
        <meta name="viewport" content="width=device-width, initial-scale=1">
        <script>
        function oAuth2() {
            var URL = 'https://notify-bot.line.me/oauth/authorize?';
            URL += 'response_type=code';
            URL += '&client_id={{.ClientID}}';
            URL += '&redirect_uri={{.CallbackURL}}';
            URL += '&scope=notify';
            URL += '&state=NO_STATE';
            window.location.href = URL;
        }
    </script>
    </head>
    <body>
        <button onclick="oAuth2();"> 連結到 LineNotify 按鈕 </button>
	</body>
`

const formTmpl = `
<!DOCTYPE html>
<html lang="tw">
    </head>
    <body>
    <form action="/connecting">
        Service Code:<br>
        <input type="text" name="code" value="123456"><br>
        <input type="text" name="token" value='{{.TOKEN}}' readonly="readonly"><br>
        <input type="submit" value="Submit">
    </form>

    </body>
`

const connectedTmpl = `
<!DOCTYPE html>
<html lang="tw">
    </head>
    <body>
    {{.TOKEN}} and {{.CODE}} is connected !! thanks.
    </body>
`
