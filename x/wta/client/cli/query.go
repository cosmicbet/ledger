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
		GetDrawCmd(),
		GetTicketsCmd(),
		GetParamsCmd(),
	)

	return cmd
}

// GetDrawCmd allows to query the details of the next draw
func GetDrawCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "draw",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.Draw(context.Background(), &types.QueryDrawRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(&res.Draw)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

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
