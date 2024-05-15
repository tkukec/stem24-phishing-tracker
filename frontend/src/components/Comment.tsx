import {IComment} from "@/interfaces/PhishingEventIntefaces";

const Comment = (props: IComment) => {
    const {comment, createdAt, username} = props
    return (
        <div className={"flex-col"}>
            <div className={"flex justify-between"}>
                <div>{username}</div>
                <div>{createdAt.toLocaleDateString()}</div>
            </div>
            <div>{comment}</div>
        </div>
    );
};

export default Comment;
