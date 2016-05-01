# Why?

@alindeman nerd sniped me into writing a tool that'd proxy zipball downloads into tarballs on the fly. This came about mostly because not all operating systems supporting piping STDOUT into `unzip`, but you can always do it with `tar`.

I would not recommend this in production.