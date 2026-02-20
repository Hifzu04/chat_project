import { MessageSquare } from "lucide-react";

const HomeSpace = () => {

  return (
    //flex-1 tells this component: "Take up all the remaining space that the sidebar isn't using." 
    // Without this, the component might shrink to only be as wide as the text inside it.
    <div className="w-full flex flex-1 flex-col items-center justify-center p-16 bg-base-100/50">
      {/* //space-y-6: Magic utility. Instead of writing margin-bottom on the icon, the heading, and the paragraph, 
    this class automatically inserts vertical space (1.5rem) between every child element inside this div. */}
      <div className="max-w-md text-center space-y-6">
        {/* Icon Display */}
        <div className="flex justify-center gap-4 mb-4">
          <div className="relative">
            <div
              className="w-16 h-16 rounded-2xl bg-primary/10 flex items-center
             justify-center animate-bounce"
            >
              <MessageSquare className="w-8 h-8 text-primary " />
            </div>
          </div>
        </div>

        {/* Welcome Text */}
        <h2 className="text-2xl font-bold">Welcome to ChatNest!</h2>
        <p className="text-base-content/60">
          Select a friend from the sidebar to start chatting
        </p>
      </div>
    </div>
  );
};

export default HomeSpace;





{/*

The CSS Concept: The "Anchor" (position: relative)
In CSS, if you want to freely drag an item around and pin it to a specific corner (using absolute), you have to give it a boundary.

relative (The Parent): Acts as the Anchor or the boundary box.

absolute (The Child): Acts as the Floating Item. It looks up the HTML tree until it finds a parent with relative, and locks itself to that box.

Why is it in your code if it's not being used?
It is almost certainly leftover code or future-proofing.

If you remember your Sidebar.js file, you had this exact code to show the green "Online" dot:


<div className="relative">
  <img src={user.profile_pic} />
  <span className="absolute bottom-0 right-0 bg-green-500" />
</div>

And in your MessageInput.js file, you used it to put the "X" button on the image preview:
JavaScript
<div className="relative">
  <img src={imagePreview} />
  <button className="absolute -top-1.5 -right-1.5"> <X /> </button>
</div>
In HomeSpace.js, whoever wrote the code probably wrapped the bouncing icon in <div className="relative"> because they were thinking: 
"Maybe later I will want to put a little notification badge, a sparkle icon, or a green dot next to this bouncing message square."
 By putting relative there now, the parent is already turned into an "Anchor." If you ever decide to add an absolute badge inside that div later,
  it will perfectly attach to the icon instead of flying off to the corner of the entire web page!
 */}