package main

var hs256 = "EyllPgDqUmu9T+ununAWNL02fKXjQfo+QWQNpqDU6TA="
var rs256 = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEApuiq3ip9hokrXtxPKEaN
JXmBqZuPt0pNINdGe8f87Yc0lgr4QMwOEQSfRKOnZZFGMdqEvIWR/8fTiOVW2oMr
MeudgxmdHECHQmLRAxacNTjLT8gKdhIgAzmqKSSLg4pDchg2J7M7T4KHODtcEY+t
MY9bpi2CCRndnoPp6ieUFRk+eEaABEdTb/4tFqHEkg4QCv/UUoWEgiKCpL02AqE9
c1iDf6KRgeQFEJUQGCu+RTiCqbIel8mTQoNY9zS/A4pPZ+7fsNEhFfF8FzcbuUd+
FezxxjscLyDwvo2892A0Vh8F/Yf5z/hgRXiPUu9yycwSM01MlknU6SAoWumrl3VZ
BwIDAQAB
-----END PUBLIC KEY-----`
var es256 = `-----BEGIN PUBLIC KEY-----
MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE3jFoUF28hp03kO7mOxyS6Y/FYo4K
zkULzGv3f8S9riY3EeUkoTKSUFR7Q9RClea/dgS9pKCFqU6UkFjfAeqx1Q==
-----END PUBLIC KEY-----`

var jwtHs256User1 = "eyJraWQiOiIxIiwiYWxnIjoiSFMyNTYifQ.eyJzdWIiOiJ1c2VyOjEifQ.hUsPgxd0QTaaH5BhkcYVdSxmJXATA5v_KAWeYAL6uVM"
var jwtHs256User2 = "eyJraWQiOiIxIiwiYWxnIjoiSFMyNTYifQ.eyJzdWIiOiJ1c2VyOjIifQ.qyrfIx0E6194v5XIhk22_geV0aSgoEdxVk4dv_FOiEo"
var jwtRs256User3 = "eyJraWQiOiIyIiwiYWxnIjoiUlMyNTYifQ.eyJzdWIiOiJ1c2VyOjMifQ.GN_5KmuZ-hKs9JXpI_TQnxvCCFRBgvc71vpBlDvgkOaXP5hALm1hfmv0ZAknLfFgKi-XqU5tYOwGWJmEJHrhi4fMusugpyqBKLNdEfZA1meZ3AYlBCCPPoS0B6i9hAPhdiiMu7L18i_0l_oQJkkQ9Dn7RT8ts7kb2M2JkR0WQMmxb8oOM_xNiji2AVk5x45DY4JI_4AWx7aoHOQDb4M35BpHiDA9qxiubHOaEIjaHSgyhrf_61cidCuPsGftbNQijh0qf6yBp10598bVewSsjH9uazOedOa7j7MPWz3X8e_0HphJZMxpdS0gVs_IczabOFvBNDJLcHdi89NhebZqnA"
var jwtRs256User4 = "eyJraWQiOiIyIiwiYWxnIjoiUlMyNTYifQ.eyJzdWIiOiJ1c2VyOjQifQ.bt9ap6JFlg9uLWVSv2Mp-Y3edLTiDrAjyrVqMWQW5gOObZEBCpbgw1VLS2yHmE3VvkOs-O1oCJ8u_w6LB0OXr4k_YOHSbJiIep7a2fm09UDwSTt6jZpOqxtO8DkGvbF5563jDFfGqIUQ81YbVPVZziKqGLpSXwyP6MfPJwVOqoMgMma0K-cJs72dST4y3cO441PYzC-Xrf7yiGe87QV9XmZRvxKjyZUmMhUupz21pLVYOyrZ0w2B7J8a4vU_9joiydEpHEoimLPC6qd2-56DHqY-gw4Q_O4GPYwoiOR6it4R6zcAlT6EUG19uejCWZP2qkm_CLTt8DV-km9KZ-NHKw"
var jwtEs256User5 = "eyJraWQiOiIzIiwiYWxnIjoiRVMyNTYifQ.eyJzdWIiOiJ1c2VyOjUifQ.eL5porEZxBCCAdzljkZkvAMkNb5ZiYMDHNfl2WKgAyRs7QDT8GmTHkUrAnUmTrwAAs4y91Z30oZqL600KceL_Q"
var jwtEs256User6 = "eyJraWQiOiIzIiwiYWxnIjoiRVMyNTYifQ.eyJzdWIiOiJ1c2VyOjYifQ.sob2SDlJbpmf0Dt6IE3YQvkhlVK-vjTCqcKZyiIkDnt6DBcRyVP4BeRkSWlV_8VyHiQddcED-F-uZhejAR3ZFg"
var jwtHs256User7Expired = "eyJraWQiOiIxIiwiYWxnIjoiSFMyNTYifQ.eyJzdWIiOiJ1c2VyOjciLCJleHAiOjF9.jNEUTOPdeqR1i_rw6YMi0YFgecclXsCHP9_G5XwzkkQ"
var jwtHs256User8Invalid = "eyJraWQiOiIxIiwiYWxnIjoiSFMyNTYifQ.eyJzdWIiOiJ1c2VyOjgifQ.Tx0KvCj3mvxCyWcBYPF9skutisKTGK5ezMtnOmmpzUc"
