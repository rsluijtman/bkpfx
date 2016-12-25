package main

import (
  "bufio"
  "fmt"
  "io/ioutil"
  "os/exec"
  "regexp"
  "sort"
  "syscall"
)


// var rxvar = regexp.MustCompile( `\$[0-9a-zA-Z_]+` ) 
var rxvar = regexp.MustCompile( `\$\w+` ) 
var inited bool = false

func substitutevars( val *string, cfgm *map[string]string ) {
  if ( inited == false ) { 
    inited = true
  }
  var loc []int 
  loc = rxvar.FindIndex( []byte( *val ) ) 
  if ( loc != nil ) {
    var before []byte
    var after []byte
    if ( loc[0] > 0 ) {
      before = val[ : loc[0]-1 ]
    }
    if ( loc[1] < val.len ) {
      after=val[ loc[1]+1 : ]
    }
    var tmp string
    var newval string

    tmp = "found var: "+ *val + "-> "
    newval = string( before ) +
             cfgm[ string( val[ loc[0]+1:loc[1] ] ) ] +
             string( after )
    *val = tmp + newval
  }
  // return val 
}

var cfgdir = "/etc/postfix"
func main(){

  e := syscall.Chdir( cfgdir )
  if  e != nil  {
    fmt.Printf( "chdir %s failed: %s\n", cfgdir, e.Error() )
    syscall.Exit( 1 )
  }
  // copy files: /etc/postfix directory + possibly /etc/aliases

  // start postconf and read output
  var configkey []string
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
      configkey = append( configkey, string(match[1]) ) 
    }
  }
  fmt.Println( "------------------" ) 
  for key, value := range configmap {
    fmt.Println("Key:", key, "Value:", value)
  }
  fmt.Println( "------------------" ) 
  sort.Strings( configkey )
  for _, key := range( configkey ) {
    fmt.Println( "key:", key, "value:", configmap[key] )
  }
  fmt.Println( "------------------" ) 

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

  fmt.Println( "------------------" ) 
  // susbstitute $config_directory
  for key, value := range configmap {
    substitutevars( &value )
    fmt.Println("Key:", key, "Value:", value)
  }
}
