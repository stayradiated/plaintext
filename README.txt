plaintext
=========

> minimal cli tool for making plain text files accessible online

what does it do?
----------------

currently it just copies all *.txt files to *.html files

but html isn't plain text!
--------------------------

I know... but how else are you going to get links?

how do I use it?
----------------

  $ cd ~/notes
  $ ls
  monday.txt tuesday.txt
  $ ./plaintext
  $ ls
  monday.txt monday.html tuesday.txt tuesday.html

can I customise the HTML?
-------------------------

of course.

this is the default template:

  <!doctype html>
  <head><title>{{.Title}}</title></head><body><pre><code>
  {{.Content}}
  </code></pre></body>

save it as a text file and tell plaintext where to find it

  $ ls
  template.html note.txt
  $ plaintext ./template.html
  $ ls
  template.html note.txt note.html

