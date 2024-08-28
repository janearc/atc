### atc

I don't remember what ATC originally stood for, it was just an idea that I had.

In principle, ATC is an intermediary between Google, Strava, and OpenAI's LLM service.
Additionally, ATC alters the Strava RPE data based upon heart rate streams to conform
to the ["TSS"](https://help.trainingpeaks.com/hc/en-us/articles/204071944-Training-Stress-Scores-TSS-Explained) metrics
(so, TSS, TSB, ATL, and CTL). This is primarily for purposes of estimating performance
training load, as well as for determining what training to undertake in order to meet
a given athletic goal.

The envisioned application flow for ATC is:
* Pull an athlete's Strava activities
* Normalize activities to TSS
* Determine current level of fitness
* Solicit goal performances or fitness level from the athlete
* Construct workouts based upon past performance and estimates which should aid the athlete in achieving these goals

Importantly, TSS is useful for the sports it was designed to quantify: triathlon, or,
swimming, biking, and running. This means that non-triathlon sports are excluded from the
activities polled from Strava, and ATC should not be considered helpful for non-triathlon
sports (although for "just" swimming, biking, or running, the logic should work fine).

### build and run atc

ATC has a specific Strava application id, and a secret key for that application id. Accordingly, it is
unlikely that someone without those secrets would be able to run this application locally. This being
said, there are two config files.

#### `config/config.yml`

```yaml
server:
  port: <port to listen on>
  redirect_uri: <the redirect uri you want to use for the oauth flow>

strava:
  url: "https://www.strava.com"

athlete:
  run:
    threshold_hr: <threshold for running, ex: 171>
  swim:
    threshold_hr: <threshold for swimming, ex: 144>
  bike:
    threshold_hr: <threshold for cycling, ex: 164>
```

#### `config/secrets.yml`

```yaml
strava:
  client_id: "124662"
  client_secret: "your strava app secret"
openai:
  api_key: "your openai API access key"
```

I think that having a client id should be all you need to authenticate to strava. As of currently,
28 Aug 2024, there is no functional openai integration. As soon as I get that corrected, ATC will have
full support for the training related things it is designed for. I think if you had your own openai key
that would also work for you, because the logic (the structured queries sent to openai) are in the code
itself and not particular to my access key.

---
author: jane mf arc, jane.arc@pobox.com
license: i do not consider this released software at the moment and i would appreciate you contact me before using it.
