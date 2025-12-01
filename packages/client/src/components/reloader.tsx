"use client"

import { FC, useEffect } from "react"
import { useRouter } from "next/navigation"

interface ReloaderProps {
    interval?: number
}

export const Reloader: FC<ReloaderProps> = ({ interval }) => {
    const router = useRouter()

    useEffect(() => {
        const reload = () => {
            router.refresh()

            setTimeout(reload, interval || 1000)
        }
        reload()
    }, [])
  return null
}