# ping2music

## What is ping2music?

ping2music is a project that I had and wanted to see if I could do it.

For now it's in a really early draft status.

### Principle:

The app is a server that listens to ICMP requests and generate a music video
based on that.

Step 1: Take a key and a duration.

Step 2: For that duration, for each ping:
 * take a random note in the key
 * take a random note duration
 * display a random color for a random duration on the screen

Step 3: At the end of the duration, take a new key and a new song duration

### Random ideas:

* Live feed on FB
* Each video is saved at the end of the duration and published on FB
* Use the source IP and the timestamp to generate a seed the seed
* Imitate different instruments?
* integrate different effects for a given instrument?
* What about rhythm instruments? Do we use drum loops or rely on the fact that
  pings are generally sent on regular intervals?
