# Gopher Castle Adventure

Your task is to write a simple text adventure game in Go. The player should be able to enter commands to move around, look at objects, and pick up objects. The game should feature multiple locations, and at least one puzzle which requires an object to solve it.

For example, a typical play session might look something like this:

```
./adventure

Welcome to Gopher Castle! While on a walking holiday, you notice an
interesting-looking castle in the distance, and decide to investigate.

Outside the Castle

You're standing outside a tumbledown castle, in front of a massive and
heavily-scarred wooden front door to the north. It is standing
invitingly half-open. Do you dare to enter?

What would you like to do?

> ENTER CASTLE

Castle Keep

You are in the castle's very solid-looking stone keep. A massive
wooden door leads out to the south, and a passageway leads north,
further into the castle.

> S

You make a move toward the door, but an evil-looking giant crab
appears from behind a secret panel to block your path!

> FIGHT CRAB

That doesn't sound like a good idea. The crab waves its massive claws
in your direction.

> N

Castle Courtyard

You are in a spacious courtyard inside the castle walls. There is a
Rust book here.

> EXAMINE BOOK

"Writing games in Rust". This looks interesting.

> TAKE BOOK

Taken.

> GO SOUTH

Castle Keep

You are in the castle's very solid-looking stone keep. A massive
wooden door leads out to the south, and a passageway leads north,
further into the castle. A giant crab blocks the door to the outside.

> GIVE BOOK TO CRAB

"Hey, that looks cool!" The crab sits down with the book and starts
reading it avidly. It seems to be paying you no further attention.

> OUT

You slip past the distracted crab and make it to freedom!

Well done, you escaped from Gopher Castle with your life! Would you
like to play again (Y/N)?

> N

Thanks for playing Gopher Castle!

```


Source: https://github.com/bitfield/adventure