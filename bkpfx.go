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


var rxvar = regexp.MustCompile( `\$[0-9a-zA-Z_]+` ) 
// var rxvar = regexp.MustCompile( `\$\w+` ) 
var notfirst bool;

func substitutevars( val *string, cfgm map[string]string ) {
  var loc []int 
  var bval []byte ;

  if notfirst==false {
    notfirst=true
    for k,v := range( cfgm ) {
      println( k,v, cfgm[k] )
    }
    println( "------------------" )
  }
  bval=[]byte(*val)
  loc = rxvar.FindStringIndex( *val ) 
  if ( loc != nil ) {
println("val :",val,*val)
println( "match loc[0]: ", loc[0], "loc[1]: ", loc[1] )
    var before []byte
    var after []byte
    if ( loc[0] > 0 ) {
      before = bval[ : loc[0] ]
    }
    if ( loc[1] < len( bval ) ) {
      after=bval[ loc[1] : ]
    }
    var newval string
    var bkey []byte
    bkey = bval[ loc[0]+1:loc[1] ] 
    
    var key string
    // key has/needs trailing space !!
    // key=string(bkey ) + " "
    key=string(bkey )
    // tmp = "found var: "+ *val + "-> "
    if v,e := cfgm[key];e {
      println( "replacing key: ",key,"with value:",v )
    } else {
      println( "cfgm",key,"is not defined",e )
      return
    }

    newval = string( before ) + cfgm[ key ] + string( after )
    // newval = string( before ) + v + string( after )
    // *val = tmp + newval
    *val = newval
  }
  // return val 
}

var cfgdir = "/etc/postfix"
func main(){
  var debug int
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
  rx := regexp.MustCompile( `^(.*\S)\s*=\s*(.*)$` )
  debug=9
  for ; e == nil ; line, _, e = r.ReadLine() {
   
    if debug > 10 {
      fmt.Println( string(line) )
    }
    match = rx.FindSubmatch( line )
    if match == nil {
      fmt.Printf( "no match: %s\n", string(line) )
    } else {
/*      fmt.Printf( "key: %s - value: %s\n", string( match[1] ),
                   string( match[2]) ) */
      configmap[ string( match[1] ) ]  = string( match[2] )
      configkey = append( configkey, string(match[1]) ) 
    }
  }
  fmt.Println( "------------------" ) 
/*
  for key, value := range configmap {
    fmt.Println("Key:", key, "Value:", value)
  }

  fmt.Println( "------------------" ) */
  if debug > 10 {
    sort.Strings( configkey )
    for _, key := range( configkey ) {
      fmt.Println( "key:", key, "value:", configmap[key] )
    }
    fmt.Println( "------------------" ) 
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

  fmt.Println( "------------------" ) 
  // susbstitute $config_directory
  var key string
  var value string
  for key, value = range configmap {
    substitutevars( &value, configmap )
    fmt.Println("Key:", key, "Value:", value)
    println( key,configmap[key])
  }
}
