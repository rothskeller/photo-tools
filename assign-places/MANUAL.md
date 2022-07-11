# assign-places

The assign-places tool is a streamlined way of assigning place tags to pictures
from an event based on their geocodes.

    usage: assign-places file...

assign-places opens a browser window showing the first image, and a map with a
marker where the first image is geocoded.  It then asks for the place tag for
that image, which can be entered using the special formats below.  Once a place
tag has been entered, assign-places moves on to the next image, and repeats
until all images have been handled.  (Images that don't contain any geocoding
are skipped.)

When entering the place tag for an image, the default is the (first) place tag
it already has, if any; otherwise, the place tag that was assigned to the
previous image.  The following things can be entered:

- A blank line.  This accepts the default.
- A full path starting with a slash.  This is a complete replacement of the
  default with a different string.  (The leading slash is omitted from the
  actual tag.)
- A relative path.  This is added to the default and then simplified, and the
  result used as the new tag.  For example, entering "../foo" would replace the
  last component of the default with "foo".

Note that assign-places does not assign multiple place tags to the same image.
If an image already has multiple place tags, only the first one is changed.
