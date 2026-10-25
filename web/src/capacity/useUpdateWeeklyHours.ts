import { useMutation, useQueryClient, type InfiniteData } from '@tanstack/react-query'
import { updateWeeklyHours, type UpdatedPerson } from './api'
import type { CapacityPage } from './types'

export function useUpdateWeeklyHours(onSuccess: () => void) {
  const queryClient = useQueryClient()

  return useMutation({
    mutationKey: ['updateWeeklyHours'],
    mutationFn: ({ personId, weeklyHours }: { personId: number; weeklyHours: number }) =>
      updateWeeklyHours(personId, weeklyHours),
    onSuccess: (person) => {
      queryClient.setQueriesData<InfiniteData<CapacityPage>>(
        { queryKey: ['capacity'] },
        (data) => replacePersonInCapacityData(data, person),
      )
      void queryClient.invalidateQueries({ queryKey: ['capacity'] })
      onSuccess()
    },
  })
}

export function replacePersonInCapacityData(
  data: InfiniteData<CapacityPage> | undefined,
  updated: UpdatedPerson,
): InfiniteData<CapacityPage> | undefined {
  if (!data) return data
  return {
    ...data,
    pages: data.pages.map((page) => ({
      ...page,
      people: page.people.map((person) =>
        person.id === updated.id
          ? { ...person, weeklyHours: updated.weeklyHours, etag: updated.etag }
          : person,
      ),
    })),
  }
}
