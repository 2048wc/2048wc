Each game and each move should be kept in a redis noSQL database.

The back-end needs to support:
A) Requirements for all games
  1) Showing all games between nth best through kth best game.
  2) Showing all games between nth most recently played through kth most recently played.
  3) Showing all games with scores between m through p.
B) Requirements for all games within a challenge
  1), 2) and 3) analogous to A)
C) Requirements for all games played by a single player
  1), 2) and 3) analogous to A)

For the constants n, k, m, p you should assume, depending on the set of requirements A), B) or C):
1) n, k to be any 2 integers, such that 0 < n < k <= z, where z is either:
  A) a total number of games played by all players
  B) a total number of games played within a given challenge
  C) a total number of games played by a given player.
2) m, p to be any 2 integers such that z1 <= m <= p <= z2, where z1 and z2 are either:
  A) respectively the lowest and the highest score ever achieved by any player.
  B) respectively the lowest and the highest score by all the players attempting a given challenge.
  C) respectively the lowest and the highest score achieved by a given person.
All ranges are left- and right- inclusive.

To cope with the amount of data, a lazy iterator should be implemented.
