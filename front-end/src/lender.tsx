
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

  
