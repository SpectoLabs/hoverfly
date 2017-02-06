hoverctl start
hoverctl destination "^.*api.*com" --dry-run https://api.github.com
hoverctl destination "^.*api.*com" --dry-run https://api.slack.com
hoverctl destination "^.*api.*com" --dry-run https://github.com
hoverctl stop