/*
Package msgid allows for creating a unique key traceable to a particular msgid
generator or instance.  The msgid is a composite of a time (to millisecond
resolution), a spawn (or node) id and an arbitrary sequence/counter value.


Usage

Create a MsgIdGenerator with a unique (to that instance)  spawner id.  Spawner
Ids would have to be coordinated across multiple instances. Millis are stored
in an 64 bit integer. A 32 bit spawn id and sequence are stored in 64 bit
integer.

For convenience "NextKey" will return a base64 encoded string that can
be decoded to a MsgId instance.

*/
package msgid
