/* eslint-disable */
export type Maybe<T> = T | null;
export type InputMaybe<T> = Maybe<T>;
export type Exact<T extends { [key: string]: unknown }> = { [K in keyof T]: T[K] };
export type MakeOptional<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]?: Maybe<T[SubKey]> };
export type MakeMaybe<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]: Maybe<T[SubKey]> };
export type MakeEmpty<T extends { [key: string]: unknown }, K extends keyof T> = { [_ in K]?: never };
export type Incremental<T> = T | { [P in keyof T]?: P extends ' $fragmentName' | '__typename' ? T[P] : never };
/** All built-in and custom scalars, mapped to their actual values */
export type Scalars = {
  ID: { input: string; output: string; }
  String: { input: string; output: string; }
  Boolean: { input: boolean; output: boolean; }
  Int: { input: number; output: number; }
  Float: { input: number; output: number; }
};

export type Book = {
  __typename?: 'Book';
  added?: Maybe<Scalars['String']['output']>;
  author: Scalars['String']['output'];
  available: Scalars['Boolean']['output'];
  edition: Scalars['String']['output'];
  id: Scalars['String']['output'];
  ownerId: Scalars['String']['output'];
  title: Scalars['String']['output'];
  updated?: Maybe<Scalars['String']['output']>;
};

export type BookList = {
  __typename?: 'BookList';
  data?: Maybe<Array<Maybe<Book>>>;
  next: Scalars['String']['output'];
  prev: Scalars['String']['output'];
  total: Scalars['String']['output'];
};

export type Query = {
  __typename?: 'Query';
  bookByOwnerId: Array<Book>;
  books: BookList;
  userByEmail?: Maybe<User>;
  userById?: Maybe<User>;
  users: Array<User>;
};


export type QueryBookByOwnerIdArgs = {
  id: Scalars['String']['input'];
};


export type QueryBooksArgs = {
  limit: Scalars['String']['input'];
  offset: Scalars['String']['input'];
};


export type QueryUserByEmailArgs = {
  email: Scalars['String']['input'];
};


export type QueryUserByIdArgs = {
  id: Scalars['String']['input'];
};

export type User = {
  __typename?: 'User';
  active: Scalars['Boolean']['output'];
  books?: Maybe<Array<Maybe<Book>>>;
  email: Scalars['String']['output'];
  firstName: Scalars['String']['output'];
  id: Scalars['String']['output'];
  lastName: Scalars['String']['output'];
  version: Scalars['Int']['output'];
};

/** All built-in and custom scalars, mapped to their actual values */
export type Scalars = {
  ID: { input: string; output: string; }
  String: { input: string; output: string; }
  Boolean: { input: boolean; output: boolean; }
  Int: { input: number; output: number; }
  Float: { input: number; output: number; }
};

export type Book = {
  __typename?: 'Book';
  added?: Maybe<Scalars['String']['output']>;
  author: Scalars['String']['output'];
  available: Scalars['Boolean']['output'];
  edition: Scalars['String']['output'];
  id: Scalars['String']['output'];
  ownerId: Scalars['String']['output'];
  title: Scalars['String']['output'];
  updated?: Maybe<Scalars['String']['output']>;
};

export type BookList = {
  __typename?: 'BookList';
  data?: Maybe<Array<Maybe<Book>>>;
  next: Scalars['String']['output'];
  prev: Scalars['String']['output'];
  total: Scalars['String']['output'];
};

export type Query = {
  __typename?: 'Query';
  bookByOwnerId: Array<Book>;
  books: BookList;
  userByEmail?: Maybe<User>;
  userById?: Maybe<User>;
  users: Array<User>;
};


export type QueryBookByOwnerIdArgs = {
  id: Scalars['String']['input'];
};


export type QueryBooksArgs = {
  limit: Scalars['String']['input'];
  offset: Scalars['String']['input'];
};


export type QueryUserByEmailArgs = {
  email: Scalars['String']['input'];
};


export type QueryUserByIdArgs = {
  id: Scalars['String']['input'];
};

export type User = {
  __typename?: 'User';
  active: Scalars['Boolean']['output'];
  books?: Maybe<Array<Maybe<Book>>>;
  email: Scalars['String']['output'];
  firstName: Scalars['String']['output'];
  id: Scalars['String']['output'];
  lastName: Scalars['String']['output'];
  version: Scalars['Int']['output'];
};
