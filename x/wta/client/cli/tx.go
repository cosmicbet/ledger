package cli

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"

	"github.com/cosmicbet/ledger/x/wta/types"
)

// NewTxCmd returns a root CLI command handler for all x/staking transaction commands.
func NewTxCmd() *cobra.Command {
	stakingTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "WTA (Winner-Take-All) transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	stakingTxCmd.AddCommand(
		NewBuyTicketsCmd(),
	)

	return stakingTxCmd
}

// NewBuyTicketsCmd returns the Cobra command allowing to buy a specific amount of tickets for the next draw
func NewBuyTicketsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "buy-tickets [quantity]",
		Short: "Buy the specified amount of tickets for the next draw",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			quantity, err := strconv.Atoi(args[0])
			if err != nil {
				return err
			}

			msg := types.NewMsgBuyTickets(int32(quantity), clientCtx.GetFromAddress().String())
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
