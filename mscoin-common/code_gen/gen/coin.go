package gen


type Coin struct {
    Id  int  `json:"id" from:"id"`
    Name  string  `json:"name" from:"name"`
    CanAutoWithdraw  int  `json:"canAutoWithdraw" from:"canAutoWithdraw"`
    CanRecharge  int  `json:"canRecharge" from:"canRecharge"`
    CanTransfer  int  `json:"canTransfer" from:"canTransfer"`
    CanWithdraw  int  `json:"canWithdraw" from:"canWithdraw"`
    CnyRate  float64  `json:"cnyRate" from:"cnyRate"`
    EnableRpc  int  `json:"enableRpc" from:"enableRpc"`
    IsPlatformCoin  int  `json:"isPlatformCoin" from:"isPlatformCoin"`
    MaxTxFee  float64  `json:"maxTxFee" from:"maxTxFee"`
    MaxWithdrawAmount  float64  `json:"maxWithdrawAmount" from:"maxWithdrawAmount"`
    MinTxFee  float64  `json:"minTxFee" from:"minTxFee"`
    MinWithdrawAmount  float64  `json:"minWithdrawAmount" from:"minWithdrawAmount"`
    NameCn  string  `json:"nameCn" from:"nameCn"`
    Sort  int  `json:"sort" from:"sort"`
    Status  int  `json:"status" from:"status"`
    Unit  string  `json:"unit" from:"unit"`
    UsdRate  float64  `json:"usdRate" from:"usdRate"`
    WithdrawThreshold  float64  `json:"withdrawThreshold" from:"withdrawThreshold"`
    HasLegal  int  `json:"hasLegal" from:"hasLegal"`
    ColdWalletAddress  string  `json:"coldWalletAddress" from:"coldWalletAddress"`
    MinerFee  float64  `json:"minerFee" from:"minerFee"`
    WithdrawScale  int  `json:"withdrawScale" from:"withdrawScale"`
    AccountType  int  `json:"accountType" from:"accountType"`
    DepositAddress  string  `json:"depositAddress" from:"depositAddress"`
    Infolink  string  `json:"infolink" from:"infolink"`
    Information  string  `json:"information" from:"information"`
    MinRechargeAmount  float64  `json:"minRechargeAmount" from:"minRechargeAmount"`
}