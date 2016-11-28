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
  var configmap map[string]string
  configmap = make( map[string]string )
  cmd := exec.Command( "/usr/sbin/postconf" )
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
      fmt.Printf( "no match: %s\n", string(line) )
    } else {
      fmt.Printf( "key: %s - value: %s\n", string( match[1] ),
                   string( match[2]) )
      configmap[ string( match[1] ) ]  = string( match[2] )
    }
  }

  for key, value := range configmap {
    fmt.Println("Key:", key, "Value:", value)
  }

  /* files do not need to be in config_directory, but usually this
     is only aliases, so only lookup alias_database, then copy files from
     config_directory
  */

  files, e := ioutil.ReadDir( cfgdir )
  if ( e != nil ) {
    fmt.Printf( "ReadDir %s failed: %s\n", cfgdir, e.Error() )
    syscall.Exit(1)
  }
  for i, file := range files {
    fmt.Printf( "%4d %s\n", i, file.Name() )
  }
}
