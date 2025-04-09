import { useQuery } from "@apollo/client"
import {loadErrorMessages} from "@apollo/client/dev"
import { gql } from "@apollo/client";


const GET_BOOK_LIST = gql(/* GraphQL */`
  query GetBookList($limit :String!,$offset :String!){
  books(limit: $limit,offset: $offset){
  data{
    id
    title
    author
  }
    prev
    next
    total
  }
  }`);



export default function GetBooks(){
        const {loading,data,error} = useQuery(
        GET_BOOK_LIST,{variables:{limit:"5",offset:"0"}});
      loadErrorMessages()
      if (error){return <h1>error</h1>};
      if (loading){return <h3>loading...</h3>}
      return(
        <div>
            <ul>
              // TODO: insert card and display
                {data.books.data.map((book :any) =>
                  <li key={book.id}>
                    title: {book.title}
                    author: {book.author}
                    </li>
                )}
            </ul>
        </div>
      )

  
};

// export default function BasicMenu() {
//   const [anchorEl, setAnchorEl] = React.useState<null | HTMLElement>(null);
//   const open = Boolean(anchorEl);
//   const handleClick = (event: React.MouseEvent<HTMLButtonElement>) => {
//     setAnchorEl(event.currentTarget);
//   };
//   const handleClose = () => {
//     setAnchorEl(null);
//   };
// const availableClick = ( event: React.MouseEvent<HTMLButtonElement>)=>{
//   console.log("event:",event.button)
//         const {loading,data} = useQuery(
//         GET_BOOK_LIST,{variables:{limit:"5",offset:"0"}});
//       return(
//         <div>
//           {loading ? (<p>loading...</p>):(
//             <p>{data}</p>
//             )}
//         </div>
// )};  

//   return (
//     <div>
//   <Container>
//     <Box sx={{ flexGrow: 1 }}>
//       <Grid2 container spacing={4}>
//         <Grid2 size={8}></Grid2>
//         <Grid2 size={4}>
//           <Item>
//             <ul>
//       <p><Button
//         id="basic-button"
//         aria-controls={open ? 'basic-menu' : undefined}
//         aria-haspopup="true"
//         aria-expanded={open ? 'true' : undefined}
//         onClick={handleClick}
//       >
//         Profile
//       </Button></p>
//       <Menu
//         id="basic-menu"
//         anchorEl={anchorEl}
//         open={open}
//         onClose={handleClose}
//         MenuListProps={{
//           'aria-labelledby': 'basic-button',
//         }}
//       >
//         <MenuItem onClick={handleClose}>My account</MenuItem>
//         <MenuItem onClick={handleClose}>Logout</MenuItem>
//       </Menu>
//     </ul>
//   </Item>
// </Grid2>
// </Grid2>
// <Stack direction="row" spacing={2}>
//       <Button
//         id="available"
//         onClick={availableClick}>Available</Button>
//       <Button id="lended">Lended</Button>
//       <Button id="dueSoon">Due soon</Button>
//       </Stack>
// </Box>
// </Container>
// </div>
//   );
// }
