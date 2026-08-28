# Structure

The files are named after the openSubsonic specifications (https://opensubsonic.netlify.app/docs/endpoints/). Then the file name. For example:
[system](https://opensubsonic.netlify.app/categories/system/) contains the endpoints getLicense, getOpenSubSonicExtensions, ping and tokenInfo.
So ping.go implements HandleGetLicense, HandleGetOpenSubsonicExtensions, HandlePing and HandleTokenInfo respectively
