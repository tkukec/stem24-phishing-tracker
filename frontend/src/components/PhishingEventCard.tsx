import {
    Card,
    CardContent,
    CardDescription,
    CardFooter,
    CardHeader,
    CardTitle,
} from "@/components/ui/card";
import { IPhishingEvent } from "@/interfaces/PhishingEventIntefaces";
import { useMemo } from "react";
import { Button } from "./ui/button";

interface PhishingEventCardProps {
    event: IPhishingEvent;
}

const TEXT_LIMIT = 70;

const PhishingEventCard = ({ event }: PhishingEventCardProps) => {
    const truncatedDescription = useMemo(() => {
        return event.description.length > TEXT_LIMIT
            ? event.description.slice(0, TEXT_LIMIT) + "..."
            : event.description;
    }, [event.description]);
    return (
        <Card className="min-w-[250px] max-w-[370px] min-h-[150px] flex flex-col">
            <CardHeader className="items-start">
                <CardTitle>{event.name}</CardTitle>
                <CardDescription className="text-left">{truncatedDescription}</CardDescription>
            </CardHeader>
            <CardContent className="text-left text-sm flex-grow">
                <div className="flex flex-wrap gap-2">
                    {event.keyword.map((keyword, index) => (
                        <span key={index} className="px-2 py-1 bg-red-600 rounded-full text-white">
                            {keyword}
                        </span>
                    ))}
                </div>
            </CardContent>
            <CardFooter className="flex justify-between">
                <p className="text-sm font-light opacity-45">{event.createdAt.toDateString()}</p>
                <Button>More info</Button>
            </CardFooter>
        </Card>
    );
};

export default PhishingEventCard;
