# build

go build -o ss .

You can move or link the executable to your PATH to have it available from everywhere

move ss /usr/local/bin  # for instance

# Usage

ss [--except=word1,regex1,..] [--no-exec] text path1 path2 ... pathN 

--no-exec  excludes executable files from searh
--except   excludes files/dirs/patterns from search, they can be regular expressions or literals or mixed

Also, you can create .ss_except file in your home directory with an entry per line, this file will be 
taken in account as if the entries will be passed with --except=...
