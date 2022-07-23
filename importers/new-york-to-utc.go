package importers

import "time"

//NewYorkToUTC converts an integer timestamp in New York to
//UTC
func NewYorkToUTC(i int64) time.Time {
	newYorkLocation, err := time.LoadLocation("America/New_York")

	if err != nil {
		panic(err)
	}

	fakeUtc := time.Unix(i, 0)

	newYork := time.Date(
		fakeUtc.Year(),
		fakeUtc.Month(),
		fakeUtc.Day(),
		fakeUtc.Hour(),
		fakeUtc.Minute(),
		fakeUtc.Second(),
		fakeUtc.Nanosecond(),
		newYorkLocation,
	)

	return newYork.UTC()
}
