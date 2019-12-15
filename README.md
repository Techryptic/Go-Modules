# Go-Modules

Recently switched over from using Ruby threaded Modules, to go-rountines. Here are some kernel tweaks to achieve the best results:

  - Check your systems: *ulimit -a*
  
  - You'll want to set your: ulimit -n 100000
  
  - https-test.txt = All the links from badssl.com for cert testing.
  
  https://medium.com/@pawilon/tuning-your-linux-kernel-and-haproxy-instance-for-high-loads-1a2105ea553e
  
  https://unix.stackexchange.com/questions/108603/do-changes-in-etc-security-limits-conf-require-a-reboot
  
  
  
Web-Module.go - Supports both HTTP/HTTPS.
