# assign-people

The assign-people tool is a streamlined way of assigning people metadata to
pictures from an event.

    usage: assign-people file...

assign-people starts by reading the people metadata from the files and assigning
an abbreviation to each person found (generally their initials, in lowercase).

Then, for each file listed, it displays the list of known abbreviations, and
marks which ones are currently tagged. Then it asks for a new list. When
asking for a new list, it accepts the following answers:

- A whitespace-separated list of abbreviations: it clears all previous people
  from the file and adds the ones identified by those abbreviations. If there
  are any unknown abbreviations on the list, it will ask for them to be defined.
- A '+' sign followed by a list of abbreviations: as above, except that it
  adds to the existing people rather than replacing them.
- A '-' sign followed by a list of abbreviations: as above, except that it
  removes the specified people from the list rather than adding them.
- The string "-ALL": it removes all people from the file.
- A blank line: it makes no changes to the file.
