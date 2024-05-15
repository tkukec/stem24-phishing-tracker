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
import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogFooter,
    DialogHeader,
    DialogTitle,
    DialogTrigger,
} from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";

const EventModal = () => {
    return (
        <Dialog>
            <DialogTrigger asChild>
                <Button variant="outline">Edit Profile</Button>
            </DialogTrigger>
            <DialogContent className="sm:max-w-[425px]">
                <DialogHeader>
                    <DialogTitle>Edit profile</DialogTitle>
                    <DialogDescription>
                        Make changes to your profile here. Click save when you're done.
                    </DialogDescription>
                </DialogHeader>
                <div className="grid gap-4 py-4">
                    <div className="grid grid-cols-4 items-center gap-4">
                        <Label htmlFor="name" className="text-right">
                            Name
                        </Label>
                        <Input id="name" value="Pedro Duarte" className="col-span-3" />
                    </div>
                    <div className="grid grid-cols-4 items-center gap-4">
                        <Label htmlFor="username" className="text-right">
                            Username
                        </Label>
                        <Input id="username" value="@peduarte" className="col-span-3" />
                    </div>
                </div>
                <DialogFooter>
                    <Button type="submit">Save changes</Button>
                </DialogFooter>
            </DialogContent>
        </Dialog>
    );
};

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
        <Dialog>
            <Card className="min-w-[250px] max-w-[370px] min-h-[150px] flex flex-col">
                <CardHeader className="items-start">
                    <CardTitle>{event.name}</CardTitle>
                    <CardDescription className="text-left">{truncatedDescription}</CardDescription>
                </CardHeader>
                <CardContent className="text-left text-sm flex-grow">
                    <div className="flex flex-wrap gap-2">
                        {event.keyword.map((keyword, index) => (
                            <Button
                                key={index}
                                className="px-2 py-1 bg-theme rounded-full text-white"
                            >
                                {keyword}
                            </Button>
                        ))}
                    </div>
                </CardContent>
                <CardFooter className="flex justify-between">
                    <p className="text-sm font-light opacity-45">
                        {event.createdAt.toDateString()}
                    </p>
                    <DialogTrigger asChild>
                        <Button>More info</Button>
                    </DialogTrigger>
                </CardFooter>
            </Card>
            <DialogContent className="sm:max-w-[425px]">
                <DialogHeader>
                    <DialogTitle>More information</DialogTitle>
                </DialogHeader>
                <div className="w-87">
                    <div className="flex mb-2">
                        <h2 className="w-35">Name</h2>
                        <h2 className="w-65">{event.name}</h2>
                    </div>
                    <div className="flex mb-2">
                        <h2 className="w-35">Description</h2>
                        <h2 className="w-65">{event.description}</h2>
                    </div>
                    <div className="flex mb-2">
                        <h2 className="w-35">Brand</h2>
                        <h2 className="w-65">{event.brand}</h2>
                    </div>
                    <div className="flex mb-2">
                        <h2 className="w-35">Malicious URL</h2>
                        <h2 className="w-65">{event.maliciousUrl}</h2>
                    </div>
                    <div className="flex mb-2">
                        <h2 className="w-35">Domain registration date</h2>
                        <h2 className="w-65">{event.domainRegistrationDate.toDateString()}</h2>
                    </div>
                    <div className="flex mb-2">
                        <h2 className="w-35">Status</h2>
                        <h2 className="w-65">{event.status}</h2>
                    </div>
                    <div className="flex mb-2">
                        <h2 className="w-35">Keywords</h2>
                        <div className="w-65">
                            {event.keyword.map((keyword, index) => (
                                <span key={index}>{keyword}, </span>
                            ))}
                        </div>
                    </div>
                    <div className="flex mb-2">
                        <h2 className="w-35">DNS records</h2>
                        <div className="w-65">
                            {event.dnsRecords.map((record, index) => (
                                <span key={index}>{record}, </span>
                            ))}
                        </div>
                    </div>
                </div>
            </DialogContent>
        </Dialog>
    );
};

export default PhishingEventCard;
