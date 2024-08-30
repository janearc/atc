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

### magical stuff

For this discussion, it's important to understand

$$
\text{CTL} = \frac{\sum \text{TSS} \text{ over 42 days}}{42}
$$

CTL is sometimes called "fitness." It is an individual, somewhat specific, metric for estimating how
difficult a workout was for an athlete _to recover from_. IF, or _intentisty factor_, is a measure of
how relatively difficult a given workout was for an athlete based upon their previous performance.

I have been using Strava, TrainingPeaks, and ChatGPT/OpenAI for planning my training and assessing my
performance for a while now. One of the things that's been tricky for me as an athlete is "given a certain
goal that I have, how do I train to achieve that goal?" This is kind of mysterious because coaches often
don't have a real good answer for how they specifically plan training, or alternatively, they have one
training program for all their athletes. Neither of these is great, because every athlete is different.
Moreover, when we have these metrics like TSS, there's no reason to give everyone the same workouts.
Because we know, subjectively, where every athlete is at, and we have some reasonable guardrails for
how much more volume or intensity (understand that TSS and CTL are essentially a "budget" you can plan
your training within, and this makes volume and intensity somewhat interchangeable in terms of planning)
can be added in a given week (such as "don't exceed more than 10% volume increase week over week"), we
have the ability to just say, "my performance right now is 10km takes one hour at `IF=1` and `TSS=100`."

Where this becomes kind of magical is when we then say, "well, what if I want to run 12km in one hour?"
Based upon an athlete's CTL (specifically, their CTL for each event), we can [SWAG](https://en.wikipedia.org/wiki/Scientific_wild-ass_guess)
what the CTL would need to be in order for a 12km, 1 hour, 100 TSS workout to be realistic. Crucially,
once we know what CTL is required to attain that performance, and we know what the athlete's capacity is
to sustain daily TSS (sometimes called TSSd), we can figure out how long it would take for them to hit
that CTL. Thus we know what \Delta \text{CTL} is, and we can break that into "chunks of TSS," which
is of course spread over 42 days, and you divide that by however many workouts per week an athlete might
have.

$$
\text{Future CTL} = \text{Current CTL} + \Delta \text{CTL}
$$

So while we might divide this by 42, if in practice an athlete has 3 runs per week, and perhaps that
long run is one of those days, and one of those days is intervals, we know several things.

* our long run is probably going to be 90 minutes at `IF=.7` (`$ \text{TSS} = \text{Duration (hours)} \times \text{IF}^2 \times 100 $`)
* our intervals run is probably going to be shorter, perhaps 45 minutes, at `IF=.85`
* this leaves us with 127 TSS for those two runs. if we want to arrive at, for example, 75 TSSd, and we assume equal commitment to swim, bike, and run, this leaves us with a budget of about 50 TSS (525 TSS per week, divided by 3 sports, minus 127 TSS for the other two runs)
* given we have a budget then of 50 TSS, and we know what the athlete's performance is, we can plan a run that is either short and intense (such as 30 minutes at `IF=1`), or a longer more relaxed run (80 minutes at `IF=.6`)
* where this gets really cool is that we can take this data, which we know about the athlete and what their delta TSS is, and we can ask openai, "hey, so I want to run this duration at this intensity" and the data that it has access to combined with the language interface allows us to do some very complicated planning that would be impossible to do in Strava, or TrainingPeaks, or in fact ChatGPT alone. Furthermore, it is difficult for _coaches with access to all this data_ to plan workouts to this degree of athlete specificity.

I'll add more to this document to reflect this as it gets written but I'm already committing this into the `planning/` directory.

---
author: jane mf arc, jane.arc@pobox.com

license: i do not consider this released software at the moment and i would appreciate you contact me before using it.



CTL is calculated as:


