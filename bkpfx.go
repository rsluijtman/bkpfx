package main

import ( 
  "fmt"
  "log"
  "os/exec"
  "syscall"
)

var cfgdir = "/etc/postfix"
func main(){

  e := syscall.Chdir( cfgdir )
  if ( e != nil ) {
    fmt.Printf( "chdir %s failed: %s\n", cfgdir, e.Error() ) ;
    syscall.Exit( 1 );
  }
  // copy files: /etc/postfix directory + possibly /etc/aliases
  // start postconf and read output
  out,e := exec.Command( "/usr/sbin/postconf", "-d" ).Output()
  if ( e != nil ) {
    log.Fatal(e)
  }
  fmt.Println(out)
}
