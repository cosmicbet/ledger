package cli

import (
	"context"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"

	"github.com/cosmicbet/ledger/x/wta/types"
)

// GetQueryCmd returns the parent command for all x/wta CLi query commands. The
// provided clientCtx should have, at a minimum, a verifier, Tendermint RPC client,
// and marshaller set.
func GetQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the wta module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		GetNextDrawCmd(),
		GetPastDrawsCmd(),
		GetTicketsCmd(),
		GetParamsCmd(),
	)

	return cmd
}

// GetNextDrawCmd allows to query the details of the next draw
func GetNextDrawCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "next-draw",
		Short: "Get the details of the next draw to be held",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.NextDraw(context.Background(), &types.QueryNextDrawRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(&res.Draw)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetPastDrawsCmd allows to query all the past draws
func GetPastDrawsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "past-draws",
		Short: "Get the details of all the past draws",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			res, err := queryClient.PastDraws(context.Background(), types.NewPastDrawsRequest(pageReq))
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(&res.Draws)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "past draws")

	return cmd
}

// GetTicketsCmd returns the Cobra command allowing to query all the sold tickets for the next draw
func GetTicketsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "tickets",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			res, err := queryClient.Tickets(cmd.Context(), types.NewTicketsRequest(pageReq))
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(&res.Tickets)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "all tickets")

	return cmd
}

// GetParamsCmd allows to query the current parameters
func GetParamsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "params",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.Params(context.Background(), &types.QueryParamsRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(&res.Params)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
