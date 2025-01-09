package commands

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/authelia/authelia/v4/internal/authentication"
)

// LoadProvidersAuthenticationRunE is a special PreRunE that loads the user authentication provider into the CmdCtx.
func (ctx *CmdCtx) LoadProvidersAuthenticationRunE(cmd *cobra.Command, args []string) (err error) {
	ctx.providers.UserProvider = getAuthenticationProvider(ctx)

	if err = doStartupCheck(ctx, providerNameUser, ctx.providers.UserProvider, false); err != nil {
		return err
	}

	return nil
}

// UserChangePasswordRunE updates user's password .
func (ctx *CmdCtx) UserChangePasswordRunE(cmd *cobra.Command, args []string) (err error) {
	if len(args) != 2 {
		return errors.New("invalid number of parameters")
	}

	var username = args[0]

	var password = args[1]

	if err := ctx.providers.UserProvider.UpdatePassword(username, password); err != nil {
		ctx.log.Fatal(err)
	}

	ctx.log.Info("password changed!")

	return nil
}

// UserShowInfoRunE shows user info.
func (ctx *CmdCtx) UserShowInfoRunE(cmd *cobra.Command, args []string) (err error) {
	var username = args[0]

	var details *authentication.UserDetails

	if details, err = ctx.providers.UserProvider.GetDetails(username); err != nil {
		ctx.log.Fatal(err)
	}

	fmt.Printf(`User '%s' Info:
	Display Name:	%s
	Email:		%s
	Groups:		%v
`, username, details.GetDisplayName(), strings.Join(details.GetEmails(), ", "), strings.Join(details.GetGroups(), ", "))

	return nil
}

// UserAddRunE adds a user.
func (ctx *CmdCtx) UserAddRunE(cmd *cobra.Command, args []string) (err error) {
	if ctx.config.AuthenticationBackend.DB == nil {
		return errors.New("this command is only available for 'db' authentication backend")
	}

	var username = args[0]

	var flags = cmd.Flags()

	password, err := flags.GetString("password")
	if err != nil {
		return err
	}

	email, err := flags.GetString("email")
	if err != nil {
		return err
	}

	displayName, err := flags.GetString("display-name")
	if err != nil {
		return err
	}

	groups, err := flags.GetStringSlice("group")
	if err != nil {
		return err
	}

	provider := ctx.providers.UserProvider.(*authentication.DBUserProvider)

	err = provider.AddUser(username, displayName, password, authentication.WithEmail(email), authentication.WithGroups(groups))
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	fmt.Println("user added.")

	return nil
}

// UserDeleteRunE deletes a user.
func (ctx *CmdCtx) UserDeleteRunE(cmd *cobra.Command, args []string) (err error) {
	if ctx.config.AuthenticationBackend.DB == nil {
		return errors.New("this command is only available for 'db' authentication backend")
	}

	var username = args[0]

	provider := ctx.providers.UserProvider.(*authentication.DBUserProvider)

	err = provider.DeleteUser(username)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	fmt.Println("user deleted.")

	return nil
}
