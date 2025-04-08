import { Container, Box } from "@mui/material";
import { styled } from '@mui/material/styles';
import Paper from '@mui/material/Paper';
import { data, NavLink } from "react-router";
import Grid from '@mui/material/Grid2';
// import Card from '@mui/material/Card';
// import CardContent from '@mui/material/CardContent';
// import CardMedia from '@mui/material/CardMedia';
// import Typography from '@mui/material/Typography';
// import CardActions from "@mui/material/CardActions";
// import Button from "@mui/material/Button";
import axios from "axios"; 
import * as gqlType from "./gql/graphql"
// import createHttpError from "http-errors";

const Item = styled(Paper)(({ theme }) => ({
  backgroundColor: '#fff',
  ...theme.typography.body2,
  padding: theme.spacing(1),
  textAlign: 'center',
  color: theme.palette.text.secondary,
  ...theme.applyStyles('dark', {
    backgroundColor: '#1A2027',
  }),
}));

const backgroundURL = "http://localhost:8080/query";


async function getBooks(){
  return axios({url:backgroundURL,
    method:'POST',
  data:{
    query: `{
          books(offset: "0", limit: "15") {
            data {
              title
              author
            }
            prev
            next
            total
          }
        }`
      }})
}

function extractJSON(r :any) {
  const jsonResult :any= JSON.stringify(r)
  return JSON.parse(jsonResult);
}

function getBookList(jsonData :any) :any{
  const bookList :any=jsonData.data.data.books.data as gqlType.BookList
  const totBooks :bigint=bookList.length
  const prev :bigint=  jsonData.data.data.books.prev
  const next :bigint=  jsonData.data.data.books.next
  const total :bigint = jsonData.data.data.books.total
  return {
    books: bookList,
    len: totBooks,
    prev: prev,
    next: next,
    total:total,
  }
}

export default function Lender() {
    getBooks().then(response=>{
    const jsonValue = extractJSON(response)
    const books = getBookList(jsonValue)
      console.log("data list:",books)
        })
      // const response = await getBooks()
    //   response.catch(error=>error)
    //   const dataList = response.then(data=>{
    //     const jsonData = extractJSON(data)
    //     const dataValues = getBookList(jsonData)
    //     console.log(dataValues.books)
    //           })
    // console.log("datalist:",dataList)
  return (
    <div className="lender">
      <Container>
        <Box sx={{ flexGrow: 1 }}>
          <Grid container spacing={4}>
            <Grid size={8}></Grid>
            <Grid size={4}>
              <Item>
                <ul>
                  <p> profile: first.last@email.com</p>
                  <NavLink to={"/logout"}> logout</NavLink>
                </ul>
              </Item>
            </Grid>
          </Grid>
        </Box>
      </Container>
    <div className = "bookList">
      <Container>
    <Box sx={{ flexGrow: 1 }}>
    </Box>
    </Container>
      </div>
    </div >
  );
};

