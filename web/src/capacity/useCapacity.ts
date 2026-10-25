import { useInfiniteQuery } from '@tanstack/react-query'
import { fetchCapacityPage } from './api'

export function useCapacity(from: string, to: string) {
  return useInfiniteQuery({
    queryKey: ['capacity', from, to],
    initialPageParam: undefined as string | undefined,
    queryFn: ({ pageParam, signal }) =>
      fetchCapacityPage({ from, to, cursor: pageParam, signal }),
    getNextPageParam: (lastPage) => lastPage.pagination.nextCursor,
    staleTime: 30_000,
    retry: 1,
  })
}
