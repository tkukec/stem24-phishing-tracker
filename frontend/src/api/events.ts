import useAxiosPrivate from "@/hooks/useAxiosPrivate";
import { PhishingEventSearchData } from "@/interfaces/PhishingEventIntefaces";
import { useInfiniteQuery } from "@tanstack/react-query";

export const EventsService = {
    useGetEvents: (data: PhishingEventSearchData) => {
        const axios = useAxiosPrivate();

        return useInfiniteQuery({
            queryKey: [
                "all-events",
                data.name,
                data.startDate,
                data.endDate,
                data.brand,
                data.domainName,
                data.keywords,
            ],
            queryFn: async ({ pageParam = 0 }) => {
                const res = await axios.get(`/events`, {
                    params: {
                        page: pageParam,
                        ...data,
                    },
                });

                return res.data;
            },
            initialPageParam: 0,
            getNextPageParam: (lastPage) => (lastPage.hasNextPage ? lastPage.nextPage : undefined),
        });
    },
};
