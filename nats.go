package main

import(
  "os"
  "log"
  "flag"
  "strings"
  "runtime"

  "github.com/apcera/nats"
  "github.com/codegangsta/cli"
  )

  var message = "Usage: nats sub/nats pub [-s Server] [--ssl] [-t] <subject>"
  var index = 0

func usage() {
  log.Fatalf(message)
}

func printMsg(m *nats.Msg, i int){
  index += 1
  log.Printf("[#%d] Received on [%s]: '%s'\n", i, m.Subject, string(m.Data))
}

func main(){
  app := cli.NewApp()
  app.Name = "nats"
  app.Usage = "Nats Pub and Sub - Go Client"
  app.Action = func(C *cli.Context) {
    println(message)
  }
  app.Commands = []cli.Command{
    {
      Name:       "pub",
      ShortName:  "p",
      Usage:      message,
      Action: func(c *cli.Context){
        var urls = flag.String("s", nats.DefaultURL, "The nats server URLs (separated by comma)")
        var ssl = flag.Bool("ssl", false, "Use Secure Connection")

        log.SetFlags(0)
        flag.Usage = usage
        flag.Parse()

        args := flag.Args()
        if len(args) < 1 {
          usage()
        }

        opts := nats.DefaultOptions
        opts.Servers = strings.Split(*urls, ",")
        for i, s := range opts.Servers {
          opts.Servers[i] = strings.Trim(s, " ")
        }

        opts.Secure = *ssl

        nc, err := opts.Connect()
        if err != nil {
          log.Fatalf("Can't connect: %v\n", err)
        }

        subj, msg := args[0], []byte(args[1])

        nc.Publish(subj, msg)
        nc.Close()

        log.Printf("Published [%s] : '%s'\n", subj, msg)
      },
    },
    {
      Name:       "sub",
      ShortName:  "c",
      Usage:      message,
      Action:  func(c *cli.Context){
        var urls = flag.String("s", nats.DefaultURL, "The nats server URLs (separated by comma)")
        var showTime = flag.Bool("t", false, "Display timestamps")
        var ssl = flag.Bool("ssl", false, "Use Secure Connection")

        log.SetFlags(0)
        flag.Usage = usage
        flag.Parse()

        args := flag.Args()
        if len(args) < 1 {
          usage()
        }

        opts := nats.DefaultOptions
        opts.Servers = strings.Split(*urls, ",")
        for i, s := range opts.Servers {
          opts.Servers[i] = strings.Trim(s, " ")
        }
        opts.Secure = *ssl

        nc, err := opts.Connect()
        if err != nil {
          log.Fatalf("Can't connect: %v\n", err)
        }

        subj, i := args[0], 0

        nc.Subscribe(subj, func(msg *nats.Msg) {
          i += 1
          printMsg(msg, i)
          })

          log.Printf("Listening on [%s]\n", subj)
          if *showTime {
            log.SetFlags(log.LstdFlags)
          }

          runtime.Goexit()
      },
    },
  }

  app.Run(os.Args)
}
