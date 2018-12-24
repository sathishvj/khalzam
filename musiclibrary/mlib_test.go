package musiclibrary

import (
	_ "github.com/lib/pq"
	"sync"
	"testing"
)

func TestIndexing(t *testing.T) {
	musicLib, err := Open()
	checkErr(err)
	wg := sync.WaitGroup{}
	defer musicLib.Close()
	wg.Add(1)
	go func() {
		defer wg.Done()
		musicLib.DeleteSong("Modjo - Lady (Hear Me Tonight)")
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		musicLib.DeleteSong("Beastie Boys - Intergalactic")
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		musicLib.DeleteSong("Mogwai - Travel Is Dangerous")
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		musicLib.DeleteSong("journeydontstop")
	}()
	wg.Wait()

	wg.Add(1)
	go func() {
		defer wg.Done()
		musicLib.InsertSong("../resources/Modjo - Lady (Hear Me Tonight).mp3")
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		musicLib.InsertSong("../resources/Beastie Boys - Intergalactic.mp3")
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		musicLib.InsertSong("../resources/Mogwai - Travel Is Dangerous.mp3")
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		musicLib.InsertSong("../resources/journeydontstop.mp3")
	}()
	wg.Wait()
}

func TestRecogniseOnly(t *testing.T) {
	musicLib, err := Open()
	checkErr(err)
	defer musicLib.Close()
	wg := sync.WaitGroup{}
	defer musicLib.Close()
	wg.Add(1)
	go func() {
		defer wg.Done()
		musicLib.RecogniseSong("../resources/intergalactic_sample.mp3")
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		musicLib.RecogniseSong("../resources/travel_chorus_sample.mp3")
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		musicLib.RecogniseSong("../resources/travel_sample.mp3")
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		musicLib.RecogniseSong("../resources/modjo_lady_sample.mp3")
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		musicLib.RecogniseSong("../resources/intergalacticnew.mp3")
	}()
	wg.Wait()
}
