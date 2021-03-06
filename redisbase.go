package main

import (
"fmt"
"github.com/fzzy/radix/redis"
"time"
"github.com/go-martini/martini"
"github.com/martini-contrib/render"
"github.com/codegangsta/martini-contrib/binding"
"strconv"
"math/rand"
"strings"
)

//Fields for the base
type Item struct {
	Band string `form:"Band"`
	Album string `form:"Album"`
	Members string `form:"Members"`
	Year int `form:"Year"`
}

type Finder struct {
	Search string `form:"Finder"`
}


func GetBands(r* redis.Client, target string)([]string, bool){
	results, error:= r.Cmd("keys", target).List()
	if error != nil {
		return nil, false
	}

	result := make([]string, len(results))
  	for i, artist := range results {
  	 	result[i] = r.Cmd("hmget", artist, "Album").String()
  	 }
  	return result, true
}

func Init(){
	r, err := InitRedis()
	defer closeRedis(r)
	rand.Seed( time.Now().UTC().UnixNano())
	localstore := make([]Item,10)
	fmt.Println(localstore)
	if err == nil{
		InitMartini(r)
	}
}
func InitMartini(r* redis.Client){
	 m := martini.Classic()
	 m.Use(render.Renderer())

  	 m.Get("/", func(ren render.Render){
    	ren.HTML(200,"index", nil)
  	 })

  	 m.Post("/", binding.Form(Item{}), func(item Item, ren render.Render){
  	 		newalbum := newAlbum(item.Band, item.Album, item.Members, item.Year)
  	 		r.Cmd("hmset", strings.ToLower(item.Band) + ":" + strconv.Itoa(rand.Int()), newalbum)
  	 		ren.HTML(200, "index", nil)
  	 	})

  	 m.Get("/find", func(ren render.Render){
  	 		ren.HTML(200, "find", nil)
  	 	})

  	 //Get list of all bands
  	 m.Get("/bands", func(fnd Finder, ren render.Render){
  	 	results, _:= r.Cmd("hgetall", fnd.Search).Hash()
  	 	fmt.Println(results)
  	 	})

  	 //Find bands by name, album or year
  	 m.Post("/find", binding.Form(Finder{}), func(fnd Finder, ren render.Render){
  	 		if(len(fnd.Search) == 0){
  	 			ren.HTML(200, "find", nil)
  	 		}
  	 		bands, ok := GetBands(r, "*" + strings.ToLower(fnd.Search) + "*")
  	 		if ok {
  	 			ren.HTML(200, "find", map[string][]string {"results": bands})
  	 		}
  	 	})

  	 m.NotFound(func(ren render.Render){
  	 		ren.HTML(200, "notfound", nil)
  	 	})
  	 m.Run()
}

func InitRedis() (*redis.Client, error){
	return redis.DialTimeout("tcp", "127.0.0.1:6379", time.Duration(10)*time.Second)
}

func closeRedis(r* redis.Client){
	defer r.Close()
}

func newAlbum(band string, album string, members string, year int) map [string] string{
	return map[string]string {
		"Band": strings.ToLower(band),
		"Album": strings.ToLower(album),
		"Members": strings.ToLower(members) /*strings.Split(members, ",")*/,
		"Year": strings.ToLower(strconv.Itoa(year)),
	}
}
func getNames(item Item){
	fmt.Println(item.Band)
}
func main() {
    Init()
}