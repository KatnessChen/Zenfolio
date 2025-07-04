import type { TRADE_TYPE } from '@/constants'

export type TradeType = (typeof TRADE_TYPE)[keyof typeof TRADE_TYPE]
