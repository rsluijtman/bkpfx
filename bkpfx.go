package main

import ( 
  "fmt"
  "io/ioutil"
  "log"
  "os/exec"
  "syscall"
)

var cfgdir = "/etc/postfix"
func main(){

  e := syscall.Chdir( cfgdir )
  if  e != nil  {
    fmt.Printf( "chdir %s failed: %s\n", cfgdir, e.Error() )
    syscall.Exit( 1 )
  }
  // copy files: /etc/postfix directory + possibly /etc/aliases
  // start postconf and read output
  out,e := exec.Command( "/usr/sbin/postconf", "-n" ).Output()
  if ( e != nil ) {
    log.Fatal(e)
  }
  fmt.Print(string(out))
  
  files, e := ioutil.ReadDir( cfgdir )
  if ( e!= nil ) {
    fmt.Printf( "ReadDir %s failed: %s\n", cfgdir, e.Error() )
    syscall.Exit(1)
  }
  for i, file := range files {
    fmt.Printf( "%4d %s\n", i, file.Name() )
  }
}
