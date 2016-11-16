package main

import ( 
  "bufio"
  "fmt"
  "io/ioutil"
  "os/exec"
  "regexp"
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

  cmd := exec.Command( "/usr/sbin/postconf", "-n" )
  stdout, e := cmd.StdoutPipe()
  cmd.Start()
  r := bufio.NewReader(stdout)
  var line []byte ;
  var match [][]byte
  rx := regexp.MustCompile( `^(.+)\s*=\s*(.*)$` )
  for ; e == nil ; line, _, e = r.ReadLine() {
    fmt.Println( string(line) )
    match = rx.FindSubmatch( line )
    if match == nil {
      fmt.Printf( "no match: %s", string(line) )
    } else {
      fmt.Printf( "key: %s - value: %s\n", string( match[1] ),
                   string( match[2]) ) ;
    }
  }

  files, e := ioutil.ReadDir( cfgdir )
  if ( e != nil ) {
    fmt.Printf( "ReadDir %s failed: %s\n", cfgdir, e.Error() )
    syscall.Exit(1)
  }
  for i, file := range files {
    fmt.Printf( "%4d %s\n", i, file.Name() )
  }
}
