package cli

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/bitrise-io/go-utils/cmdex"
	"github.com/bitrise-io/go-utils/colorstring"
	"github.com/bitrise-io/stepman/stepman"
	"github.com/codegangsta/cli"
)

func printFinishShare() {
	fmt.Println()
	log.Info(" * " + colorstring.Green("[OK] ") + "Yeah!! You rock!!")
	fmt.Println()
	fmt.Println("   " + GuideTextForFinish())
	fmt.Println()
	msg := `   You can create a pull request in your forked StepLib repository,
   if you used the main StepLib repository then your repository's url looks like: ` + `
   ` + colorstring.Green("https://github.com/[your-username]/bitrise-steplib") + `

   On GitHub you can find a ` + colorstring.Green("'Compare & pull request'") + ` button, in the ` + colorstring.Green("'Your recently pushed branches:'") + ` section,
   which will bring you to the 'Open a pull request' page, where you can review and create your Pull Request.
	`
	fmt.Println(msg)
}

func finish(c *cli.Context) {
	share, err := ReadShareSteplibFromFile()
	if err != nil {
		log.Error(err)
		log.Fatal("You have to start sharing with `stepman share start`, or you can read instructions with `stepman share`")
	}

	route, found := stepman.ReadRoute(share.Collection)
	if !found {
		log.Fatalln("No route found for collectionURI (%s)", share.Collection)
	}

	collectionDir := stepman.GetCollectionBaseDirPath(route)
	if err := cmdex.GitCheckIsNoChanges(collectionDir); err == nil {
		log.Warn("No git changes!")
		printFinishShare()
		return
	}

	stepDirInSteplib := stepman.GetStepCollectionDirPath(route, share.StepID, share.StepTag)
	stepYMLPathInSteplib := stepDirInSteplib + "/step.yml"
	log.Info("New step.yml:", stepYMLPathInSteplib)
	if err := cmdex.GitAddFile(collectionDir, stepYMLPathInSteplib); err != nil {
		log.Fatal(err)
	}

	log.Info("Do commit")
	msg := share.StepID + " " + share.StepTag
	if err := cmdex.GitCommit(collectionDir, msg); err != nil {
		log.Fatal(err)
	}

	log.Info("Pushing to your fork: ", share.Collection)
	if err := cmdex.GitPushToOrigin(collectionDir, share.StepID); err != nil {
		log.Fatal(err)
	}
	printFinishShare()
}
