package main

import (
	"fmt"
	"os"
	"bufio"
	"../sort"
)

type Location struct {
	x int
	y int
}

type Land struct {
	barren bool
	color int
	location Location
}

type BarrenSection struct {
	lowerLeft Location
	upperRight Location
}

type FertileSection struct {
	color int
	lands []Land
}

var AllLand [400][600]Land
var FertileSections []FertileSection
var BarrenSections []BarrenSection

func main() {
	// Initialize lists
	AllLand = [400][600]Land{}
	FertileSections = []FertileSection{}
	BarrenSections = []BarrenSection{}

	initAllLand()

	//input := "0 292 399 307"

	getBarrenLandsFromUser()

	FindFertileSections()

	fertileSectionSizes := []int{}
	for _, section := range FertileSections {
		fertileSectionSizes = append(fertileSectionSizes, len(section.lands))
	}

	fertileSectionSizes = sort.TopDownMergeSort(fertileSectionSizes)
	fmt.Println(fertileSectionSizes)
	//fmt.Println("Number of fertile sections", len(FertileSections))
	//spew.Dump(FertileSections)
}

// Set coordinates for all land locations
func initAllLand() {
	for i, landColumn := range AllLand {
		for j := range landColumn {
			AllLand[i][j].location = Location{i, j}
		}
	}
}

func getBarrenLandsFromUser() {
	// Read from stdin
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter barren area: ")
	input, _ := reader.ReadString('\n')
	//fmt.Println(input)

	if input != "\n" {
		// Label barren section as barren
		newBarrenSection := CreateBarrenSection(input)
		BarrenSections = append(BarrenSections, newBarrenSection)
		MarkBarrenSection(newBarrenSection)
		getBarrenLandsFromUser()
	}


}

func CreateBarrenSection(input string) (barrenSection BarrenSection){
	// Add a new Barren Section to BarrenSections
	coordinates := make([]int, 4)
	_, err := fmt.Sscan(input, &coordinates[0], &coordinates[1], &coordinates[2], &coordinates[3])
	if err != nil {
		fmt.Println(err)
		return
	}

	barrenSection = BarrenSection {
		 Location {
			coordinates[0],
			coordinates[1],
		},
		Location {
			coordinates[2],
			coordinates[3],
		},
	}

	BarrenSections = append(BarrenSections, barrenSection)
	return barrenSection
}

func MarkBarrenSection(barrenSection BarrenSection) {
	for i := barrenSection.lowerLeft.x; i <= barrenSection.upperRight.x; i++ {
		for j := barrenSection.lowerLeft.y; j <= barrenSection.upperRight.y; j++ {
			AllLand[i][j].barren = true
			AllLand[i][j].color = -1
			//AllLand[i][j].location = Location{i, j}
		}
	}
}

// Recursive FloodFill causes a stack overflow :(
func FloodFillRecursive(land Land, section FertileSection) {
	if land.barren { return }
	if land.color == section.color { return }
	land.color = section.color
	section.lands = append(section.lands, land)

	FloodFillRecursive(getNeighboringLand(land, 1, 0), section)
	FloodFillRecursive(getNeighboringLand(land, -1, 0), section)
	FloodFillRecursive(getNeighboringLand(land, 0, 1), section)
	FloodFillRecursive(getNeighboringLand(land, 0, -1), section)

}

func FindFertileSections(){
	color := 1
	for i := range AllLand {
		landColumn := AllLand[i]
		for j := range landColumn {
			land := AllLand[i][j]
			// Find an uncolored node
			if !land.barren && land.color == 0 {
				// Create a new fertile section
				fertileSection := FertileSection{color, []Land{}}
				// Flood Fill new fertile section starting with the uncolored node
				FloodFill(land, &fertileSection)

				FertileSections = append(FertileSections, fertileSection)
				color++	// Set up new color for next section
			}
		}
	}
}

func FloodFill(land Land, section *FertileSection) {
	// Return if land is barren or colored - TODO: redundant?
	if land.barren {
		fmt.Println("Land is barren")
		return }
	if land.color == section.color {
		fmt.Println("Land color matches section color")
		return }

	// Color this new section of land and add it to our fertile section
	AllLand[land.location.x][land.location.y].color = section.color
	section.lands = append(section.lands, AllLand[land.location.x][land.location.y])

	queue := make([]Land, 0)
	queue = append(queue, AllLand[land.location.x][land.location.y])

	// As long as the queue is not empty, keep flood filling
	for len(queue) != 0 {
		// Pop from queue
		x := queue[0]
		queue = queue[1:]

		// Add neighboring nodes to queue
		neighbors := []Land{
			getNeighboringLand(x, 1, 0),
			getNeighboringLand(x, -1, 0),
			getNeighboringLand(x, 0, 1),
			getNeighboringLand(x, 0, -1),
		}
		//spew.Dump(neighbors)

		for _, neighbor := range neighbors {
			if neighbor.color == 0 {
				// Color the node, add it to our section, and append to queue
				AllLand[neighbor.location.x][neighbor.location.y].color = section.color
				section.lands = append(section.lands, AllLand[neighbor.location.x][neighbor.location.y])
				queue = append(queue, AllLand[neighbor.location.x][neighbor.location.y])
			}
		}

		// Remove - kick us out if we exceed maximum possible section size
		if len(queue) > 240000 {
			//spew.Dump(queue)
			panic(fmt.Sprintf("Queue is larger than maximum possible section size. Fix your algorithm"))
		}
	}

}

func getNeighboringLand(land Land, xOffset int, yOffset int) Land{
	invalidLand := Land{true, -1, Location{-1, -1}}

	if  land.location.x + xOffset < 0 ||
		land.location.x + xOffset > 399 ||
		land.location.y + yOffset < 0 ||
		land.location.y + yOffset > 599 {
			return invalidLand
	}

	//fmt.Println(land.location.x + xOffset, land.location.y + yOffset)
	x := land.location.x + xOffset
	y := land.location.y + yOffset
	//AllLand[x][y].location = Location{x, y}
	neighbor := AllLand[x][y]

	return neighbor
}